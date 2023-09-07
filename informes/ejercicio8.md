### Implementacion concurrencia servidor

Frente a cada nuevo cliente que quiera conectarse al servidor, se spawnea un nuevo proceso que atiende al cliente. Ademas, el proceso inicial, guarda un *handle* del proceso spawneado para poder hacer join luego. 
Cuando el proceso que atiende a un cliente termina, coloca el resultado de la atencion en una cola. 

> Un resultado exitoso implica que el cliente cerro la conexion con un mensaje *End* y se enviaron apuestas al servidor.
>
> En este caso, el servidor actualiza la cantidad de clientes que restan por atender.

Antes de recibir nuevas conexiones el proceso principal lee todos los resultados que hayan quedado en la cola, haciendo join al mismo tiempo de los procesos que hayan terminado.

Para escribir en el archivo de apuestas, los distintos procesos se sincronizan a travez de un *multiprocessing.Lock*.
Cuando todas las casas de apuestas terminan de enviar la apuestas realizadas en las mismas, el proceso principal adquiere el lock al archivo y lee el archivo de apuestas y obtiene los ganadores, segregandolos por casa de apuesta.


### Graceful shutdown del servidor

El servidor implementa graceful shutdown de la siguiente manera:

1. Todos los procesos hijos definen un handler de *SIGTERM*
2. Al recibir el proceso principal una señal de tipo *SIGTERM*, reenvia la señal a todos los procesos hijos ejecutando la siguiente sentencia:

```python
    hijo.terminate()
    hijo.join()
```

3. Luego de asegurarse que todos los procesos hijos finalizaron correctamente, el proceso principal libera los recursos asociados al *socket* y a la *cola*
