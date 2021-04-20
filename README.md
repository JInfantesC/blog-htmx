# go-web-server
Programa para probar un servidor de páginas utilizando distintos paquetes y técnicas de Go.

El programa tiene un gestor automático de páginas que permite incluir archivos en la carpeta `pages` y servirlos `localhost:8080\pages\` 

## Como usar
```
# En la carpeta del programa ejecutar:
$ go run .
# o si quieres seleccionar la carpeta a servir
$ go run . --dir pages
```

## Detalles del programa
Actualmente el programa solo sirve archivos de texto plano. Los archivos `.md` los devuelve 
parseados en html.

El programa utiliza `go:embed` para incluir los archivos estáticos y plantillas en el binario.

Se hace uso de `http/template` para leer el archivo index.gohtml y enviarle el objeto directorio y 
montar el html final donde navegar por los archivos.

Se utiliza el paquete `flag` para poder indicar al ejecutable que carpeta quieres servir.

He usado de este artículo 
https://www.codemio.com/2019/01/advanced-golang-tutorial-http-middleware.html
los middleware para mostrar en la 
consola algo de login y lo he ampliado con dos funciones `HandleFunc` y `Handle` que funcionan 
como sus homólogos del paquete oficial `http` con la diferencia de que implementa los middleware.

En `pagesManager.go` he creado un único tipo `Page` que de manera recursiva con `ReadDirectories` 
gestiona la carpeta `pages` a servir. Del mismo he programado la función `HandleDirectory` para 
servir esos archivos.

## Licencias
He usado:
- https://github.com/bigskysoftware/htmx
- https://github.com/franciscop/picnic
- https://github.com/russross/blackfriday

Licencias se sirven en la carpeta `pages/licenses`

