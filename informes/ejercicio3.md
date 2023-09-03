### Implementacion

1. Se incluye dentro de la carpeta *netcat* los siguientes archivos:

    - Un Dockerfile, para buildear una imagen de alpine con netcat instalado.
    - Un script *netcat-build.sh* para buildear la imagen correspondiente.
    - Un script *netcat-run.sh* para levantar el contenedor y probar asi, el funcionamiento del servidor.

2. Se incluye en el directorio *config/netcat* archivos para configurar la execucion de netcat:

    - env.txt $\longrightarrow$ define el puerto en el cual se encuentra escuchando el servidor.
    - message_file.txt $\longrightarrow$ contiene el mensaje que sera enviado al servidor.
    - run.sh $\longrightarrow$ es el script que se corre dentro del contenedor.

### Instrucciones de uso

1. Ejecutar en primer lugar el script *netcat/netcat-build.sh* para construir la imagen.
2. Ejecutar *netcat/netcat-run.sh* para correr el contenedor y realizar la prueba del servidor.

> Si se desea modificar el mensaje enviado, se puede editar el archivo *config/netcat/message_file.txt*
>
> Si el servidor cambia el puerto en el cual escucha conexiones. Se debe modificar el archivo *config/netcat/env.txt*
>
> Ante modificaciones en los archivos de configuracion, no es necesario rebuildear la imagen, ya que estos se incluyen como parte de un volumen en el contenedor.