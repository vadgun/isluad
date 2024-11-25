var agregados = {};
var infoModal = document.getElementById('infoModal');
var commenta = document.getElementById('message-text');
var serviciosseleccionados = document.getElementById('ServiciosSeleccionados');


infoModal.addEventListener('show.bs.modal', function (event) {
    let button = event.relatedTarget; // Botón que abrió el modal
    
    // Extrae información de los atributos data-bs-*
    let title = button.getAttribute('data-bs-title');
    let content = button.getAttribute('data-bs-content');
    let value = button.getAttribute('data-bs-value');

    // Actualiza el contenido del modal
    let modalTitle = infoModal.querySelector('.modal-title');
    let modalBody = infoModal.querySelector('.modal-body p');
    const agregar = document.querySelector('.modal-footer a');
    const remover = document.querySelector('.modal-footer .btn-outline-danger');
    
    modalTitle.textContent = title;
    modalBody.textContent = content;
    
    if (agregados[value]){
        commenta.value = agregados[value].Comentario;
    }else{
        commenta.value="";
    }
    
    ActualizarDiv();
    agregar.href = "Javascript:Agregar('"+value+"');";
    remover.href = "Javascript:Remover('"+value+"');"; 
});


function Agregar(s){
    let serviciobutton = document.getElementById(s);
    serviciobutton.classList.replace("btn-outline-light", "btn-outline-success")
    c = commenta.value;
    agregados[s] = {Comentario: c};
    ActualizarDiv();
    toastAdd._config.delay=2500;
    toastAdd._element.lastElementChild.innerText = serviciobutton.textContent;
    toastAdd.show();
}

function Remover(s){      
    if (s in agregados) {delete agregados[s]
        commenta.value="";
        let serviciobutton = document.getElementById(s);
        serviciobutton.classList.replace("btn-outline-success", "btn-outline-light")
        ActualizarDiv();
        toastRemove._config.delay=2500;
        toastRemove._element.lastElementChild.innerText = serviciobutton.textContent;
        toastRemove.show()
    }
}

function ActualizarDiv(){
    if (Object.keys(agregados).length === 0 ){
        serviciosseleccionados.textContent="Ningún servicio seleccionado";
    }

    if (Object.keys(agregados).length === 1){
        serviciosseleccionados.textContent="1 servicio seleccionado";
    }

    if (Object.keys(agregados).length > 1){
        serviciosseleccionados.textContent= Object.keys(agregados).length +" servicios seleccionados";
    }
    
}