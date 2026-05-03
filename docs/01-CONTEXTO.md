# Contexto

## Qué es esto

Un **agente LLM local** (sin internet) montado en formato USB-portable, con foco en **discreción visual**:

- El proceso parece una shell de Windows en pantalla, no una aplicación de IA.
- La entrada del usuario es invisible aunque alguien mire por encima del hombro o estés compartiendo pantalla.
- La salida por defecto es silenciosa: solo el siguiente prompt indica que el agente terminó.
- No genera ficheros de log ni historial persistente.
- Toda la comunicación entre el harness y el modelo va por loopback (`127.0.0.1`); no toca la red.
- Funciona offline en cualquier PC Windows 8 GB sin GPU; cero instalación, todo desde el USB.

El objetivo práctico es poder usar IA local en una pantalla compartida (Zoom, Teams, screen-share) o con alguien al lado, sin que se vea lo que escribes ni los pasos del agente.

## Qué resuelve

Un agente IA típico (Claude Code, Cursor, Aider) imprime cada acción que hace, ocupa una pestaña distinguible de IA, y depende de internet o de servicios externos. Eso está bien para escritorio personal, mal cuando:

- Compartes pantalla y no quieres que se vea el contenido de tus prompts.
- Estás en una máquina sin internet o con tráfico saliente bloqueado.
- Prefieres no dejar rastro en historial / logs / disco.
- Quieres llevarlo en un USB y usarlo sin instalar nada en la máquina destino.

Este harness ataca específicamente esos cuatro puntos.

## Qué NO es

- No es una alternativa completa a Claude Code / Cursor: la calidad del modelo local 1.5B–7B está por debajo de modelos cloud frontier. Para tareas serias de refactor / arquitectura, sigue ganando un servicio cloud.
- No es invisible al sistema operativo: Task Manager sigue listando los procesos, `netstat` sigue mostrando el puerto local, el filesystem sigue grabando los ficheros que el agente edita. La discreción es **frente al ojo humano que mira la pantalla**, no frente a auditoría administrativa.
- No está pensado para máquinas que no son tuyas, ni para redes corporativas con políticas que tienes que respetar. Es una herramienta personal.

## Para qué SÍ vale

- **Tu portátil personal** mientras tomas un café en una cafetería o un compañero pasa por detrás.
- **Sesiones de pair-programming en remoto** donde compartes pantalla pero no quieres exponer tus prompts.
- **Streams / grabaciones** donde el contenido de los prompts es privado.
- **PCs sin internet** donde quieres usar IA — siempre que sea tu propia máquina.

## Filosofía de diseño

1. **Sigilo visual, no sigilo operativo.** No intentamos ocultar el proceso al sistema; intentamos que mirar la pantalla no revele lo que está pasando.
2. **Cero rastro accidental.** Ni logs, ni historial, ni cache. Si el usuario quiere persistencia, la pide explícitamente con env var.
3. **USB-portable.** Un binario y unos `.gguf`; copia-pega y funciona en cualquier Windows 8 GB.
4. **Calidad razonable, no perfecta.** En 8 GB sin GPU, los tok/s y la calidad del modelo limitan lo que se puede hacer. El harness no oculta esa realidad: documenta tiempos honestos en `05-DURACIONES.md`.
5. **Permisos contenidos.** El agente solo puede tocar el directorio de trabajo (CWD) por defecto. Listas negras explícitas para `.git`, `.env`, claves SSH, comandos peligrosos.

## Lectura recomendada

- Si vienes a entender la arquitectura: `02-ARQUITECTURA.md`.
- Si quieres ver técnica a técnica cómo se logra el sigilo visual: `03-TECNICAS-SIGILO.md`.
- Si te importa la elección del modelo: `04-MODELOS-Y-RAM.md`.
- Si vas a desplegarlo en Windows: `06-EJECUCION-WINDOWS.md`.
