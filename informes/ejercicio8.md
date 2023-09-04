### Implementacion concurrencia servidor

Frente a cada nuevo cliente que quiera conectarse al servidor, se spawnea un nuevo proceso que atiende al cliente. Ademas, el proceso inicial, guarda un *handle* del proceso spawneado para poder hacer join luego. 
Cuando el proceso que atiende a un cliente termina, coloca el resultado de la atencion en una cola. 

> Un resultado exitoso implica que el cliente cerro la conexion con un mensaje *End*
>
> En este caso, el servidor actualiza la cantidad de clientes que restan por atender.

El proceso principal define un handle de *SIGCHLD*. Cuando un proceso hijo termina, se ejecuta este handler y el proceso principal extrae el resultado de la cola.

Para escribir en el archivo de apuestas, los distintos procesos se sincronizan a travez de un *multiprocessing.Lock*.
Cuando todas las casas de apuestas terminan de enviar la apuestas realizadas en las mismas, el proceso principal lee el archivo de apuestas y obtiene los ganadores, segregandolos por casa de apuesta.