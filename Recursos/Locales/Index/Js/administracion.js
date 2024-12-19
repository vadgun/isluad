// Obtener
function Obtener(info) {
  var request = $.ajax({
    url: "/services",
    method: "POST",
    data: { data: info },
    dataType: "html"
  });

  request.done(function (response) {
    $("#main-content").html(response);
    if (info === "Calendario") {
      crearCalendario();
    }
    if (info === "Cotizaciones") {
      new DataTable('#cotizaciones');
    }
    if (info === "Reservaciones") {
      new DataTable('#reservaciones');
    }
  });

  request.fail(function (textStatus) {
    alert("Request failed: " + textStatus);
  });
}

// Categoria
function Categoria(info) {
  var request = $.ajax({
    url: "/categoria",
    method: "POST",
    data: { data: info },
    dataType: "html"
  });

  request.done(function (response) {
    $("#tables").html(response);
    new DataTable('#example');
  });

  request.fail(function (textStatus) {
    alert("Request failed: " + textStatus);
  });
}

// Ver
function Ver(info) {
  var request = $.ajax({
    url: "/ver",
    method: "POST",
    data: { data: info },
    dataType: "json"
  });

  request.done(function (response) {
    crearModalVer(response)
  });

  request.fail(function (textStatus) {
    alert("Request failed: " + textStatus);
  });
}

// Eliminar
function Eliminar(info) {
  var request = $.ajax({
    url: "/eliminar",
    method: "POST",
    data: { data: info },
    dataType: "html"
  });

  request.done(function (response) {
    Swal.fire("Servicio eliminado correctamente");
  });

  request.fail(function (textStatus) {
    alert("Request failed: " + textStatus);
  });
}

function VerServicio(info) {
  var request = $.ajax({
    url: "/verServicio",
    method: "POST",
    data: { data: info },
    dataType: "html"
  });

  request.done(function (response) {
    console.log(response)
  });

  request.fail(function (textStatus) {
    alert("Request failed: " + textStatus);
  });
}

// Modal Ver Servicio

let modalServicio = document.getElementById('modalVerServicio');
let modalBodyServicio = document.querySelector("#modalVerServicio .modal-body");
modalServicio.addEventListener('hidden.bs.modal', event => {
  setFocus();
  modalBodyServicio.innerHTML = '';
});

modalServicio.addEventListener('show.bs.modal', function (event) {
  let button = event.relatedTarget; // Botón que abrió el modal

  // Extrae información de los atributos data-bs-*
  let title = button.getAttribute('data-bs-title');
  let cost = button.getAttribute('data-bs-cost');
  let description = button.getAttribute('data-bs-description');
  let category = button.getAttribute('data-bs-category');
  let active = button.getAttribute('data-bs-active');
  let id = button.getAttribute('data-bs-id');

  const categorias = ["Planificación y Coordinación", "Decoración y Ambientación", "Música y Entretenimiento", "Alimentos y Bebidas", "Documentos y Ceremonias", "Fotografia y Recuerdos", "Servicios Especiales"]

  // console.log(title, cost, description, category, active)

  let div = document.createElement('div');
  div.classList.add("col-12");
  let form = document.createElement('form');
  form.id = 'formUpdateService';

  cadenadehtml = `<div class="mb-3">
  <label for="titulo" class="form-label">Nombre:</label>
  <input type="text" class="form-control" id="titulo" name="titulo" value="${title}">
  </div>
  <div class="mb-3">
  <label for="descripcion" class="form-label">Descripción:</label>
  <textarea class="form-control" id="descripcion" name="descripcion" rows="5">${description}</textarea>
  </div>
  <div class="mb-3">
  <label for="costo" class="form-label">Costo:</label>
  <input type="number" class="form-control" id="costo" name="costo" value="${cost}">
  </div>
  <div class="form-check form-switch">
  <label class="form-check-label" for="activo">Activo</label>`
  if (active === "true") {
    cadenadehtml += `<input class="form-check-input" type="checkbox" role="switch" id="activo" name="activo" checked>`
  } else {
    cadenadehtml += `<input class="form-check-input" type="checkbox" role="switch" id="activo" name="activo" unchecked>`
  }
  cadenadehtml += `</div>
  <label for="categoria" class="form-label">Categoría:</label>
  <select class="form-select" id="categoria" name="categoria" aria-label="category select">`

  for (const element of categorias) {
    if (element === category) {
      cadenadehtml += `<option selected value="${element}">${element}</option>`
    } else {
      cadenadehtml += `<option value="${element}">${element}</option>`
    }
  }

  cadenadehtml += `</select>`
  form.innerHTML = cadenadehtml;
  let submitButton = document.createElement('div');
  submitButton.classList.add("container", "col-6");
  submitButton.innerHTML = `<br><div class="row">
   <input type="hidden" id="currentid" name="currentid" value="${id}">
  <button type="submit" class="btn btn-lg btn-outline-success ">Guardar servicio</button></div>`;
  form.appendChild(submitButton);
  div.appendChild(form);
  modalBodyServicio.appendChild(div);

  $('#formUpdateService').on('submit', function (e) {
    e.preventDefault(); // Prevenir el comportamiento por defecto
    // Obtener datos del formulario
    var formData = $(this).serializeArray();
    var formObject = {};
    $.each(formData, function (i, field) {
      if (field.name == "activo") {
        formObject[field.name] = Boolean(field.value);
      } else {
        formObject[field.name] = field.value;
      }
    });

    var requestData = { ...formObject };
    // console.log("en json", JSON.stringify(requestData));
    requestData["costo"] = parseFloat(requestData["costo"]);

    $.ajax({
      url: '/editarServicio',
      type: 'POST',
      contentType: 'application/json',
      data: JSON.stringify(requestData),
      success: function (response) {
        Swal.fire(response.mensaje);
        $('#modalVerServicio').modal('hide');
        Categoria(category);
      },
      error: function (error) {
        alert('Error al editar el servicio');
        console.error(error);
      }
    });
  });
});


// Crear Modal Ver
function crearModalVer(response) {
  const modalVer = new bootstrap.Modal(document.getElementById("modalVer"));
  const modalBody = document.querySelector("#modalVer .modal-body");

  const nuevafecha = new Date(response.fechaReserva);
  const opciones = {
    weekday: 'long', // Nombre del día
    day: 'numeric', // Día del mes
    month: 'long', // Nombre completo del mes
    year: 'numeric' // Año completo
  };
  const formatoLegible = new Intl.DateTimeFormat('es-ES', opciones).format(nuevafecha);
  const costo = response.costoTotal
  const formatoMex = costo.toLocaleString('es-MX', {
    style: 'currency',
    currency: 'MXN',
  });

  var disponibilidad = obtenerDisponibilidad(response.fecha);
  async function obtenerDisponibilidad(fecha) {
    try {
      const response = await $.ajax({
        url: "/date",
        method: "POST",
        data: { data: fecha },
        dataType: "json"
      });

      return response.mensaje === "Si";
    } catch (error) {
      console.error("Error en la solicitud:", error);
      return false; // Devuelve un valor por defecto en caso de error
    }
  }

  // Usar la función
  (async function () {
    const disponibilidad = await obtenerDisponibilidad(response.fecha);
  })();

  let columnas = document.createElement('div');
  columnas.classList.add('container')

  if (disponibilidad) {
    columnas.innerHTML = `
    <div class="row">
      <div class="col-6">
      <p>Creada por: <strong>${response.nombre}</strong> <br>Día: <strong>${formatoLegible}</strong><br>Costo estimado: <strong>${formatoMex}</strong><br>Telefono: <strong>${response.telefono}</strong><br>Correo: <strong>${response.correo}</strong><br>Invitados : <strong>${response.invitados}</strong></p>
      </div>
      <div class="col-6 text-center">
      <p><strong>La fecha esta disponible</strong><br><span class="icon text-success"><i class="bi bi-emoji-sunglasses-fill bigicon"></i></span></p>
      </div>
    </div>`
  } else {
    columnas.innerHTML = `
    <div class="row">
      <div class="col-6">
      <p>Creada por: <strong>${response.nombre}</strong> <br>Día:<strong>${formatoLegible}</strong><br>Costo estimado: <strong>${formatoMex}</strong><br>Telefono: <strong>${response.telefono}</strong><br>Correo: <strong>${response.correo}</strong><br>Invitados : <strong>${response.invitados}</strong></p>
      </div>
      <div class="col-6 text-center">
      <p><strong>La fecha no esta disponible</strong><br><span class="icon text-danger"><i class="bi bi-emoji-astonished bigicon"></i></i></span></p>
      </div>
    </div>`
  }

  modalBody.appendChild(columnas);

  let div = document.createElement('div');
  div.classList.add("col-12");
  let form = document.createElement('form');
  form.id = 'formModCotizacion';

  let cacheado = false
  if (cache[response.ID] === undefined) {
    cache[response.ID] = response.agregados;
  } else {
    cacheado = true
  }

  Object.keys(response.agregados).forEach(key => {
    let value = response.agregados[key];
    agregados[key] = response.agregados[key]

    if (cacheado) {
      value.costoIndividual = cache[response.ID][key].costoIndividual;
    }

    form.innerHTML = form.innerHTML + `<div class="row input-group">
    <input type="text" value="${value.nombreCompleto}" id="${key}Service" name="${key}Service" style="width: 25%;" class="form-control" aria-label="Servicio" read-only disabled>
    <textarea class="form-control" id="${key}Text" name="${key}Text" aria-label="Comentarios" style="width: 60%;" read-only disabled>${value.comentario}</textarea>
    <input type="number" onchange="agregardor(this.id +':'+ this.value +':${response.ID}');" style="width: 15%;" value="${value.costoIndividual}" id="${key}" name="${key}" class="form-control" aria-label="Costo">
    </div>`
  });

  let submitButton = document.createElement('div');
  submitButton.classList.add("container", "col-4");
  submitButton.innerHTML = `<br><div class="row">
   <input type="hidden" id="currentid" name="currentid" value="${response.ID}">
  <button type="submit" class="btn btn-lg btn-outline-success ">Crear reservación</button></div>`;
  form.appendChild(submitButton);
  div.appendChild(form);
  modalBody.appendChild(div);
  modalVer.show();
}

var modalVer = document.getElementById('modalVer')
var modalBody = document.querySelector("#modalVer .modal-body");
modalVer.addEventListener('hidden.bs.modal', event => {
  setFocus();
  modalBody.innerHTML = '';
  agregados = {};
})
var cache = new Map() // A map
var agregados = {} // An object

// agregador
function agregardor(s) {
  ss = s.split(":")
  cache[ss[2]][ss[0]].costoIndividual = Number(ss[1]);
  agregados[ss[0]].costoIndividual = Number(ss[1]);
}

// Fechas preseleccionadas (simulando que vienen de una base de datos)
var preselectedDates = [
  // { date: "2024-12-22", label: "Mi Cumpleaños" },
  // { date: "2024-12-24", label: "Cena de Navidad" },
  // { date: "2024-12-25", label: "Noche buena" },
  // { date: "2024-12-26", label: "Cena de la empresa" },
  // { date: "2024-12-31", label: "Ultima noche del año" },
  // { date: "2025-12-22", label: "Mi Cumpleaños 2025" },
];

// Crear calendario
function crearCalendario() {
  const calendarDays = document.getElementById("calendar-days");
  const monthLabel = document.getElementById("month-label");
  const yearLabel = document.getElementById("year-label");
  const monthSelector = document.getElementById("month-selector");
  const dateInput = document.getElementById("date-input");
  const prevMonthButton = document.getElementById("prev-month");
  const nextMonthButton = document.getElementById("next-month");
  const modalEvento = new bootstrap.Modal(document.getElementById("modalEvento"));
  const modalBody = document.querySelector("#modalEvento .modal-body");

  let currentDate = new Date();

  function renderCalendar(date) {
    const year = date.getFullYear();
    const month = date.getMonth();

    // Set the header
    monthLabel.textContent = date.toLocaleString("default", { month: "long" });
    yearLabel.textContent = year;

    // Get the first day and number of days in the month
    const firstDay = new Date(year, month, 1).getDay();
    const daysInMonth = new Date(year, month + 1, 0).getDate();

    // Clear previous days
    calendarDays.innerHTML = "";

    // Add empty slots for the first row
    for (let i = 0; i < firstDay; i++) {
      const emptyDiv = document.createElement("div");
      calendarDays.appendChild(emptyDiv);
    }

    // Add days of the month
    for (let day = 1; day <= daysInMonth; day++) {
      const dayDiv = document.createElement("div");
      dayDiv.textContent = day;
      dayDiv.classList.add("day");

      const formattedDate = `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`;

      // Check if the day is preselected
      const isPreselected = preselectedDates.find((d) => d.date === formattedDate);
      if (isPreselected) {
        dayDiv.classList.add("selected");
        dayDiv.setAttribute("data-label", isPreselected.label);
      }

      // Handle day click
      dayDiv.addEventListener("click", () => {
        if (isPreselected) {
          modalBody.textContent = `${formattedDate} - ${isPreselected.label}`;
          modalEvento.show();
        } else {
          alert(`Fecha disponible: ${formattedDate}`);
        }
      });

      calendarDays.appendChild(dayDiv);
    }
  }

  // Handle month navigation
  prevMonthButton.addEventListener("click", () => {
    currentDate.setMonth(currentDate.getMonth() - 1);
    renderCalendar(currentDate);
  });

  nextMonthButton.addEventListener("click", () => {
    currentDate.setMonth(currentDate.getMonth() + 1);
    renderCalendar(currentDate);
  });

  // Handle month label click
  monthLabel.addEventListener("click", () => {
    monthSelector.style.display = "block";
  });

  // Handle month label click
  yearLabel.addEventListener("click", () => {
    monthSelector.style.display = "block";
  });

  // Handle month selection
  monthSelector.addEventListener("click", (e) => {
    if (e.target.dataset.month) {
      currentDate.setMonth(e.target.dataset.month);
      renderCalendar(currentDate);
      monthSelector.style.display = "none";
    }
  });

  // Handle date input change
  dateInput.addEventListener("change", (e) => {
    const [year, month, day] = e.target.value.split("-");
    const selectedDate = `${year}-${String(month).padStart(2, "0")}-${String(day).padStart(2, "0")}`;

    const isPreselected = preselectedDates.find((d) => d.date === selectedDate);

    if (isPreselected) {
      modalBody.textContent = `${selectedDate} - ${isPreselected.label}`;
      modalEvento.show();
    } else {
      alert(`Fecha disponible: ${selectedDate}`);
    }

    currentDate = new Date(year, month - 1, day);
    renderCalendar(currentDate);
  });

  // Close month selector on outside click
  document.addEventListener("click", (e) => {
    if (!e.target.closest(".calendar-header")) {
      monthSelector.style.display = "none";
    }
  });

  // Initialize calendar
  renderCalendar(currentDate);
}

modalVer.addEventListener('show.bs.modal', event => {
  $('#formModCotizacion').on('submit', function (e) {
    e.preventDefault(); // Prevenir el comportamiento por defecto
    // Obtener datos del formulario
    var formData = $(this).serializeArray();
    var formObject = {};
    $.each(formData, function (i, field) {
      if (field.name === "currentid") {
        formObject[field.name] = field.value;
      }
    });

    var servicios = { agregados };
    var requestData = { ...formObject, ...servicios };

    console.log("en json", JSON.stringify(requestData))


    $.ajax({
      url: '/guardarreservacion',
      type: 'POST',
      contentType: 'application/json',
      data: JSON.stringify(requestData),
      success: function (response) {
        Swal.fire("Reservacion creada");
        // toastCotizacion._config.delay=2500;
        // toastCotizacion._element.lastElementChild.innerText = serviciobutton.textContent;
        // toastCotizacion.show()
        // alert(response.message)
        // location.replace("/")
        // Obtener('Cotizaciones');
      },
      error: function (error) {
        alert('Error al enviar la cotización');
        console.error(error);
      }
    });


  });
});

function setFocus() {
  console.log("focuseando");
  const elements = document.querySelectorAll('[id^="dt-search-"]');
  elements.forEach(el => el.focus());
  console.log(elements)
}

$(document).ready(function () {

  console.log(preselectedDates);

});
