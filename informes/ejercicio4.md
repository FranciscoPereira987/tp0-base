### Implementacion

#### Server

El servidor implementa graceful shutdown seteando un handler para la señal *SIGTERM*. Posteriormente se cierran todos los file descriptors.

#### Client

El cliente implementa graceful shutdown creando un channel para escuchar por la señal *SIGTERM*. Luego de cada ciclo, el cliente trata de recuperar la señal desde el mismo channel. Al momento de recibir una señal del mismo channel. El cliente procede a cerrar el channel.

### Instrucciones de uso

Para que alguno de los procesos haga un shutdown se puede invocar el siguiente comando:

```bash
docker kill --signal="SIGTERM" <containerName>
```

donde container name es el nombre del contenedor al cual queremos enviarle la señal *SIGTERM*. Este comando funcionara unicamente si el contenedor se encuentra en ejecucion actualmente. Por lo que antes debera correrse el comando:

```bash
make docker-compose-up
```

#### Ejemplos

```bash
docker kill --signal="SIGTERM" server
```

```bash
docker kill --signal="SIGTERM" client1
```
