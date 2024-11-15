var agregados = {};
var infoModal = document.getElementById('infoModal');
var commenta = document.getElementById('message-text');


infoModal.addEventListener('show.bs.modal', function (event) {
    let button = event.relatedTarget; // Botón que abrió el modal
    
    // Extrae información de los atributos data-bs-*
    let title = button.getAttribute('data-bs-title');
    let content = button.getAttribute('data-bs-content');
    let value = button.getAttribute('data-bs-value');

    // Actualiza el contenido del modal
    let modalTitle = infoModal.querySelector('.modal-title');
    let modalBody = infoModal.querySelector('.modal-body p');
    // var modalTextArea = infoModal.querySelector('.modal-body input')
    const agregar = document.querySelector('.modal-footer a');
    const remover = document.querySelector('.modal-footer .btn-outline-danger');
    
    modalTitle.textContent = title;
    modalBody.textContent = content;
    
    if (agregados[value]){
        commenta.value = agregados[value].Comentario
    }else{
        commenta.value="";
    }
    
    agregar.href = "Javascript:Agregar('"+value+"');"
    remover.href = "Javascript:Remover('"+value+"');"; 
});


function Agregar(s){
    let serviciobutton = document.getElementById(s);
    serviciobutton.classList.replace("btn-outline-light", "btn-outline-success")
    c = commenta.value;
    agregados[s] = {Comentario: c};
}

function Remover(s){      
    if (s in agregados) {delete agregados[s]}
    commenta.value="";
    let serviciobutton = document.getElementById(s);
    serviciobutton.classList.replace("btn-outline-success", "btn-outline-light")
}