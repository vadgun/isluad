var addToast = document.getElementById('liveToast')
var deleteToast = document.getElementById('deleteToast')
var cotizacionToast = document.getElementById('cotizacionToast')
var toastAdd = bootstrap.Toast.getOrCreateInstance(addToast)
var toastRemove = bootstrap.Toast.getOrCreateInstance(deleteToast)
var toastCotizacion = bootstrap.Toast.getOrCreateInstance(cotizacionToast)

$(document).ready(function () {
    $('#formCotizacion').on('submit', function (e) {
        e.preventDefault(); // Prevenir el comportamiento por defecto

        // Obtener datos del formulario
        var formData = $(this).serializeArray();
        var formObject = {};
        $.each(formData, function (i, field) {
            formObject[field.name] = field.value;
        });

        var servicios = {agregados};

        // Combinar los datos del formulario con los datos extra
        var requestData = { ...formObject, ...servicios };
        
        console.log("en json", JSON.stringify(requestData));

        $.ajax({
            url: '/guardarcotizacion',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(requestData),
            success: function (response) {
                console.log(response);
                toastCotizacion._config.delay=2500;
                // toastCotizacion._element.lastElementChild.innerText = serviciobutton.textContent;
                // toastCotizacion.show()
                alert(response.message)
                location.replace("/")
            },
            error: function (error) {
                alert('Error al enviar la cotización');
                console.error(error);
            }
        });
    });
});