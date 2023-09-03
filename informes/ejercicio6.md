### Instrucciones de uso

Para que el cliente se ejecute correctamente, es necesario tener en cuenta los siguientes parametros:

- en *config/client/config.yaml*

    - path.dataset $\longrightarrow$ indica el directorio en el cual se encuentran los archivos con las apuestas de la agencia.
    - path.file $\longrightarrow$ indica el nombre del archivo, en forma generica, en el cual se encuentran las apuestas correspondientes a la agencia.

- en *scripts/docker-compose-template.yaml.jinja*

    - Los clientes tienen configurado un volumen para los archivos de apuestas. Por ende, ante cambios en el directorio y/o cambios dentro de la imagen es necesario modificar este campo.

> Si se desea cambiar el directorio dentro del contenedor donde se encuentran las apuestas, hay que modificar tanto *config/client/config.yaml* como *scripts/docker-compose-template.yaml.jinja*

