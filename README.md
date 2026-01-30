# BadgerStatus

Monitor de estado del sistema para Badger 2040.

## Uso con Docker Compose (Recomendado para Synology NAS)

Para ejecutar este proyecto en un Synology NAS, lo más sencillo es usar Docker Compose.

1.  Copia el archivo `docker-compose.yml` a tu NAS.
2.  Asegúrate de que tu Badger 2040 esté conectado por USB. Normalmente aparecerá como `/dev/ttyACM0`.
3.  Ejecuta:

```bash
docker-compose up -d
```

### Notas para Synology:
- Se incluye `privileged: true` para asegurar el acceso al dispositivo USB.
- El volumen del sistema se mapea a `/host_volume` dentro del contenedor para que la aplicación pueda leer el uso de disco del NAS, no el del contenedor.
- Se mapean `/proc`, `/sys` y `/etc` del host para que las métricas de CPU y RAM sean las del NAS y no las del contenedor.
- Puedes cambiar el `VOLUME` en el `docker-compose.yml` si quieres monitorizar un volumen específico (ej. `/volume1`).

### Solución de problemas:
Si ves que las métricas no corresponden con las del NAS, asegúrate de reconstruir la imagen después de hacer cambios:
```bash
docker-compose up -d --build
```
Puedes ver qué configuración está cargando la aplicación con:
```bash
docker logs badger-status
```

## Variables de Entorno

| Variable | Descripción | Defecto |
| :--- | :--- | :--- |
| `VOLUME` | Punto de montaje a monitorizar | `/` |
| `SERIAL_PORT` | Puerto serie del Badger | `/dev/ttyACM0` |
| `BAUD_RATE` | Velocidad de transmisión | `115200` |
| `CRON` | Expresión cron para programar el envío (p. ej. `@every 1m` o `*/5 * * * *`) | `@every 1m` |
