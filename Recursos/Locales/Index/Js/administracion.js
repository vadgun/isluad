function Obtener(info) {
  // var menuId = $( "ul.nav" ).first().attr( "id" );
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
  });

  request.fail(function (textStatus) {
    alert("Request failed: " + textStatus);
  });
}

function Categoria(info) {
  // var menuId = $( "ul.nav" ).first().attr( "id" );
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

  // Fechas preseleccionadas (simulando que vienen de una base de datos)
  const preselectedDates = [
    { date: "2024-12-22", label: "Mi Cumpleaños" },
    { date: "2024-12-24", label: "Cena de Navidad" },
    { date: "2024-12-25", label: "Noche buena" },
    { date: "2024-12-26", label: "Cena de la empresa" },
    { date: "2024-12-31", label: "Ultima noche del año" },
    { date: "2025-12-22", label: "Mi Cumpleaños 2025" },
  ];

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