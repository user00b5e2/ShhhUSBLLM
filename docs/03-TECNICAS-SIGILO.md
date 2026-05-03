# Técnicas de discreción visual

Catálogo de los mecanismos que el harness usa para que mirar la pantalla no revele que está corriendo un agente IA. Cada técnica se explica en tres bloques:

- **Qué hace** (concepto).
- **Cómo lo implementa este proyecto** (con cita literal del código + fichero).
- **Cómo verificarlo / desactivarlo** (qué pasaría si falla y cómo lo compruebas).

---

## T1 — Nombres neutros para los binarios

### Qué hace

Los dos `.exe` del proyecto usan nombres que no llaman la atención:

- `bgupd.exe` (el harness) — sugiere "background updater".
- `hostcfg.exe` (el backend de inferencia, que internamente es `llama-server`) — sugiere "host configuration".

En Task Manager o `tasklist` aparecen mezclados con el ruido típico de procesos del sistema, en lugar de un `claude-agent.exe` o `llama-server.exe` que delaten la naturaleza de la herramienta.

### Cómo lo implementa

`build.sh:11`:
```bash
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o ../bin/bgupd.exe .
```

El llama-server descargado se renombra:
`download-llama.sh:24`: `cp "$src" bin/hostcfg.exe`

### Cómo verificarlo

`Get-Process bgupd, hostcfg` desde PS lista ambos cuando están corriendo. La técnica es solo cosmética; no oculta los procesos al sistema.

---

## T2 — Ocultación de la entrada del usuario en pantalla

### Qué hace

Lo que el usuario teclea no se renderiza, ni siquiera al pegar texto del portapapeles. Si alguien mira por encima del hombro o estás compartiendo pantalla, no ven el contenido de los prompts.

### Cómo lo implementa

**Doble capa**, para sobrevivir terminales que ignoran ANSI:

1. **Capa ANSI** — `ESC[8m` (concealed attribute). En `stealth.go:19`:
   ```go
   ansiConcealOn  = "\x1b[8m"
   ```
2. **Capa de TTY raw mode** — desactiva el echo del terminal. En `term_windows.go` y `term_unix.go`:
   ```go
   state, err := term.MakeRaw(fd)
   ```
   `MakeRaw` quita `ENABLE_ECHO_INPUT` en Windows y `ECHO` en Unix; el terminal no eco aunque pegues.

### Cómo verificarlo

Tipea con `update.bat` corriendo en la terminal. Si los caracteres aparecen, alguna de las dos capas falló. Causa típica en Windows: `conhost.exe` viejo (Windows 10 < 1809) que ignora ANSI. La capa 2 (raw mode) lo cubre.

---

## T3 — Imitación del prompt de la shell

### Qué hace

El binario imprime una cadena que **parece** el prompt real (`PS C:\Users\demo>`). Cuando el usuario abre la terminal y arranca el harness, no hay un banner "Welcome to shhh-agent" — solo lo que parece ser el prompt normal de la shell, esperando comandos.

Para alguien que mira la pantalla, no hay forma visual de distinguir "el usuario está en su PowerShell normal" de "el usuario está dentro del REPL del agente".

### Cómo lo implementa

`update.ps1` ejecuta la función real `prompt` antes de pasar el control y exporta el resultado:
```powershell
$captured = (& { prompt }) -join ''
$env:SHHH_FAKE_PROMPT = $captured.TrimEnd() + ' '
```

`stealth.go:36-50` lee la cascada:
```go
func ResolvePrompt(k ShellKind) string {
    if p := os.Getenv("SHHH_PROMPT"); p != "" { return ensureTrailingSpace(p) }
    if p := os.Getenv("SHHH_FAKE_PROMPT"); p != "" { return ensureTrailingSpace(p) }
    return FakePrompt(k)  // fallback genérico
}
```

Tres niveles, en cascada:
1. `SHHH_PROMPT` — override manual del usuario (cadena exacta).
2. `SHHH_FAKE_PROMPT` — capturado por `update.ps1` antes de arrancar.
3. Fallback — prompt genérico construido a partir de la shell detectada y el CWD.

### Cómo verificarlo

Compara el prompt que aparece tras lanzar `update.ps1` con el prompt que tenías antes. Si son distintos, la captura falló (típicamente porque el `prompt` de PS depende de variables que no se evalúan en un sub-scope; `oh-my-posh` con render asíncrono es el caso clásico). Solución: definir `SHHH_PROMPT` manualmente.

---

## T4 — Backend con ventana oculta

### Qué hace

`hostcfg.exe` (el motor LLM) no abre ventana de consola. Está cargando varios GB de modelo y usando 4 cores, pero el usuario no ve un parpadeo de ventana ni una consola adicional.

### Cómo lo implementa

`proc_windows.go:14-19`:
```go
func hideWindow(c *exec.Cmd) {
    c.SysProcAttr = &syscall.SysProcAttr{
        HideWindow:    true,
        CreationFlags: createNoWindow,  // 0x08000000 = CREATE_NO_WINDOW
    }
}
```

El `stdio` se redirige a `os.DevNull` para que tampoco escriba en la consola del padre.

### Cómo verificarlo

Lanza `update.ps1`, espera unos segundos, abre Task Manager y busca `hostcfg.exe`. Está corriendo pero sin ventana asociada. Si ves una consola fugaz parpadear al arrancar, los flags no se aplicaron (posible en builds Go muy viejos).

---

## T5 — Comunicación enteramente por loopback

### Qué hace

Toda la conversación del harness con el modelo va por `127.0.0.1:8765`. El kernel de Windows no marca eso como tráfico de red — no aparece en monitorizadores que filtran por NIC, ni dispara alertas de firewall saliente.

### Cómo lo implementa

`server.go:78`:
```go
"--host", "127.0.0.1", "--port", "8765",
```

### Cómo verificarlo

`netstat -ano | findstr 8765` muestra la conexión local con el PID de `hostcfg.exe`. `Get-NetTCPConnection -LocalPort 8765` desde PS hace lo mismo. Si `Test-NetConnection -ComputerName 127.0.0.1 -Port 8765` da `True`, está arriba.

Importante: esto significa que el harness funciona aunque el PC tenga **todo el tráfico saliente bloqueado**, porque loopback no atraviesa el firewall típico.

---

## T6 — Salida silenciosa por defecto

### Qué hace

El usuario lanza una tarea, el agente edita ficheros, no aparece nada en pantalla excepto el siguiente prompt cuando termina. El "output" del trabajo está en el filesystem, no en la consola.

Esto es fundamental para el sigilo visual durante una sesión de screen-sharing: aunque alguien esté mirando, no ve qué está haciendo el agente paso a paso.

### Cómo lo implementa

`main.go` en el bloque post-turno:
```go
if showResult { fmt.Println(strings.TrimSpace(output)) }
fmt.Print(prompt)  // sólo el prompt nuevo
```

`showResult` solo se activa con env var (`SHHH_SHOW_RESULT=1` o `SHHH_VERBOSE=1`).

### Cómo verificarlo

Lanza una tarea de creación de fichero. Tras unos segundos debe aparecer el prompt nuevo, sin output intermedio. Abre el fichero en VS Code: ahí está el resultado.

Para activar feedback: `$env:SHHH_SHOW_RESULT = "1"` antes de lanzar.

---

## T7 — Cero persistencia en disco

### Qué hace

El harness no escribe logs, no persiste el historial del REPL, no guarda cache de conversaciones. Al cerrar, lo único que queda en disco son los ficheros que el agente editó (que es el resultado deseado del trabajo).

### Cómo lo implementa

- Sin fichero de log: `cmd.Stderr = devnullW` para `hostcfg.exe`.
- Sin historial: el slice de `[]Message` del bucle ReAct vive solo en RAM del proceso.
- Sin cache: `--no-mmap` en `hostcfg.exe` (más por velocidad de USB que por sigilo, pero el efecto secundario es que no queda mmap residual).

### Cómo verificarlo

Tras una sesión, busca cualquier fichero nuevo en `%TEMP%`, `%APPDATA%`, `%LOCALAPPDATA%` con timestamp en la sesión. Solo debería aparecer `%TEMP%\hostcfg.lock` (un PID file de 3 líneas), que se borra con `update --stop`.

---

## T8 — Lockfile en %TEMP% para coordinación entre invocaciones

### Qué hace

Para que múltiples invocaciones del harness reusen el mismo backend (en lugar de cargar el modelo cada vez, lo cual son 30–90 s perdidos), el PID del backend se persiste en un fichero temporal.

### Cómo lo implementa

`server.go:39-45`:
```go
func lockPath() string {
    return filepath.Join(os.TempDir(), "hostcfg.lock")
}
```

Contiene 3 líneas: PID, ruta del modelo cargado, puerto.

### Cómo verificarlo

`Get-Content $env:TEMP\hostcfg.lock` cuando el backend está vivo. Tras `update --stop`, el fichero desaparece.

---

## Lo que NO se oculta

Honestidad sobre las cosas que **siguen siendo visibles** y no hay forma de evitar:

- **CPU al 80–100 %** durante inferencia. 4 cores ocupados es un patrón obvio en cualquier monitor.
- **RAM consumida**: 1.5–5 GB según el slot. Visible en Task Manager / `Get-Process`.
- **Puertos abiertos**: `127.0.0.1:8765` aparece en `netstat`.
- **Historial de comandos**: si alguien hace `Get-History` ve la línea `.\update.ps1`. PS persiste el historial a disco por defecto en `(Get-PSReadlineOption).HistorySavePath`. Para no persistirlo: `Set-PSReadlineOption -HistorySaveStyle SaveNothing` (manual).
- **Filesystem changes**: cada `write_file` queda registrado por el SO. Auditorías de filesystem (Sysmon, fsutil) lo capturan.

El sigilo es **frente a alguien mirando la pantalla**, no frente a alguien con privilegios administrativos en la máquina.
