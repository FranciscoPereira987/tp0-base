### Descripcion de alto nivel

El protocolo se compone de los siguientes mensajes:

1. Hello
    - Utilizado para inicializar la conexion entre cliente y servidor
2. Ack
    - Utilizado para confirmar que se recibio el mensaje anterior. No todos los mensajes se responden con un Ack
3. Err
    - Indica que ocurrio un error y que la conexion se cerro del otro lado.
4. End
    - Indica que se desea terminar la comunicacion.
5. Bet
    - Contiene los datos asociados a una apuesta
6. BetBatch
    - Es un mensaje compuesto por mensajes de tipo Bet.
7. Winners
    - Request para conocer los ganadores asociados a una agencia en particular
8. WinnersResponse
    - Respuesta a Winners con los Documentos de aquellos que ganaron el sorteo.


### Estructura de los mensajes

#### Estructura general

El header de todos los mensajes es de 4 bytes y esta compuesto de la siguiente manera:

|OP_CODE|Message Length|
|:--:|:--:|
|1 byte|3 bytes|

1. OP_CODE
    - Indica el tipo de mensaje
2. Message Length 
    - Indica la cantidad de bytes totales del mensaje.

#### Estructura por tipo de mensaje

##### Hello

Los mensajes *Hello* contienen, ademas del header, el numero de agencia del cliente que inicia la comunicacion, la estructura general se describe a continuacion.

|OP_CODE|Message length|Agency Number|
|:--:|:--:|:--:|
|0x01|8|uint32|

> Un *Hello* debe ser respondido con un *Ack*

##### Ack

Los mensajes *Ack* se componen unicamente del header, por lo que tienen la siguiente estructura:

|OP_CODE|Message Length|
|:--:|:--:|
|0x02|4|

##### Err

Los mensajes *Err*, asi como los *Ack*, estan compuestos unicamente por el header:

|OP_CODE|Message Length|
|:--:|:--:|
|0x03|4|

> Los mensajes *Err* indican que la conexion se cerro del otro extremo. Los mensajes *Err* no tienen respuesta.

##### End

Un mensaje *End* se compone unicamente por su header correspondiente:

|OP_CODE|Message Length|
|:--:|:--:|
|0xff|4|

> Los mensajes *End* se responden utilizando un *Ack* e indican que no se realizara el intercambio de mensajes posteriores.

##### Bet

El mensaje *Bet* contiene informacion sobre una apuesta. En el se incluyen los siguientes campos:

1. Name (utf-8 string)
2. Surname (utf-8 string)
3. Id (utf-8 string)
4. Birthdate (utf-8 string)
5. Beted number (uint 32, Big Endian)
6. Agency (uint 32, Big Endian)

- Los campos en formato *string* se codifican de la siguiente manera:
    - Un byte (Big Endian) que indica la longitud total del campo. (Field Length)
    - *Field Length* bytes que codifican el campo en formato utf-8

|Field Length|Field|
|:--:|:--:|
|1 byte|<256 bytes|

- Los campos *Beted number* y *Agency* se codifican como Big Endian y se componen de 4 bytes.

La estructura general de un mensaje *Bet* es la siguiente:

|OP_CODE|Message Length|
|:--:|:--:|
|0x04|Length|

con el cuerpo del mensaje compuesto por:

|||
|:--:|:--:|
|NameLength|Name|
|SurnameLength|Surname|
|IdLength|Id|
|BirthdateLength|Birthdate|
|Betted|number|
|Agency|number|

> Todo mensaje *Bet* es procesado y seguido por un *Err* en caso de error o *Ack* en caso de que se haya podido procesar la apuesta.


#### BetBatch

Los mensajes *BetBatch* son una composicion de mensajes *Bet* por lo que su estructura es la siguiente:

|OP_CODE|Message Length|
|:--:|:--:|
|0x05|Length|

con el cuerpo del mensaje compuesto por mensajes de tipo *Bet* de la siguiente manera:

|||
|:--:|:--:|
|Bet1|Bet Message|
|...|...|
|BetN|Bet Message|

> Al igual que los mensajes de tipo *Bet*, los mensajes *BetBatch* necesitan de una respuesta. O bien *Err* indicando que no se pudo procesar el batch de apuestas, o bien *Ack* indicando que se proceso correctamente.


#### Winners

Los mensajes *Winners* se componen unicamente de su header:

|OP_CODE|Message Length|
|:--:|:--:|
|0x06|4|

> Los mensajes de tipo *Winners* tienen dos respuestas posibles
>
> Una respuesta de tipo *Err* indica que el servidor todavia no realizo el sorteo.
>
> Una respuesta de tipo *WinnersResponse*, que indica quienes son los ganadores del sorteo de esa agencia.

#### Winners Response

Los mensajes de tipo Winners Response se componen del header, seguidos por la lista de los documentos de aquellas personas que hayan resultado ganadoras del sorteo. Por lo que su estructura es la siguiente:

|OP_CODE|Message Length|
|:--:|:--:|
|0x07|Length|

donde el cuerpo esta compuesto como:

|||
|:--:|:--:|
|Id1 Length|Id1|
|...|...|
|IdN Length|IdN|

Los Documentos de los ganadores, se codifican de la misma manera que los campos de los mensajes Bet. Osea que son strings utf-8 precedidas por un byte Big Endian que indica la longitud del string.


### Flujos del protocolo
