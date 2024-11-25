document.getElementById('recoverPassword').addEventListener('click', function (event) {
    event.preventDefault();
    const email = prompt('Por favor, ingresa tu correo electrónico para recuperar la contraseña:');
    if (email) {
        // Lógica para enviar correo de recuperación
        alert('Se ha enviado un correo a ' + email + ' con las instrucciones para recuperar tu contraseña.');
    }
});

// document.getElementById('loginForm').addEventListener('submit', function (event) {
//     event.preventDefault();
//     alert('Iniciar sesión presionado');
//     // Agrega tu lógica de autenticación aquí
// });

document.getElementById('togglePassword').addEventListener('click', function () {
    const passwordField = document.getElementById('password');
    const passwordFieldType = passwordField.getAttribute('type');
    if (passwordFieldType === 'password') {
        passwordField.setAttribute('type', 'text');
        this.innerHTML = '<i class="bi bi-eye-slash"></i>'; // Cambiar icono
    } else {
        passwordField.setAttribute('type', 'password');
        this.innerHTML = '<i class="bi bi-eye"></i>'; // Cambiar icono
    }
});