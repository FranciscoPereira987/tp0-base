### Modificaciones realizadas

Con el objetivo de que al realizar cambios en la configuracion no se requiera un nuevo build de las imagenes tanto del *server* como del *client* se realizaron las siguientes modificaciones:

1. Los archivos de configuracion *config.ini* y *config.yaml* se movieron a:
    - config/client/config.yaml
    - config/server/config.ini

2. En *scripts/docker-compose-template.yaml.jinja* se definieron los siguientes volumenes:

    - *./config/server:/config/* $\longrightarrow$ para el servidor.

    - *./config/client:/config* $\longrightarrow$ para cada uno de los clientes.


### Instrucciones de uso

Para modificar las configuraciones del cliente y el servidor, se deben modificar los archivos que fueron mencionados anteriormente.

- *config/server/config.ini* $\longrightarrow$ para el caso del servidor

- *config/client/config.yaml* $\longrightarrow$ para el caso de los cientes.