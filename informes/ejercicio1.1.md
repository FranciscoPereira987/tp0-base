### Instrucciones

El script se encuentra en /scripts junto con el template de jinja. Para correr el escript puede
usarse el siguiente comando:

```bash
python scripts/compose_script.py -c <clients>
```

done \<clients> es la cantidad de clientes que se quiere existan en el docker-compose.yaml

Ademas, el script se ejecuta cuando se realiza el siguiente comando:

```bash
make |CLIENTS=N| docker-compose-up
```

En caso de no indicar la cantidad de clientes, el script se ejecutara por defecto con $CLIENTS=3$

##### Ejemplos de uso

```bash
python scripts/compose_script -c 5
```

```bash
make CLIENTS=10 docker-compose-up
```