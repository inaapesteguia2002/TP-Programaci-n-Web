// URL base de tu API en Go (ajustá el puerto o la ruta si es necesario)
const API_URL = 'http://localhost:8080/partidos';

// Elementos del DOM
const formPartido = document.getElementById('form-partido');
const inputLocal = document.getElementById('equipo-local');
const inputVisitante = document.getElementById('equipo-visitante');
const listaPartidos = document.getElementById('lista-partidos');

// 1. Obtener y mostrar los partidos al cargar la página (GET)
document.addEventListener('DOMContentLoaded', obtenerPartidos);

async function obtenerPartidos() {
    try {
        const response = await fetch(API_URL);
        if (!response.ok) throw new Error('Error al obtener los partidos');
        
        const partidos = await response.json();
        renderizarPartidos(partidos);
    } catch (error) {
        console.error('Error en la petición GET:', error);
    }
}

// Renderizar los partidos dinámicamente en el <ul>
function renderizarPartidos(partidos) {
    listaPartidos.innerHTML = ''; // Limpiar la lista actual
    
    if (!partidos || partidos.length === 0) {
        listaPartidos.innerHTML = '<li>No hay partidos registrados.</li>';
        return;
    }

    partidos.forEach(partido => {
        const li = document.createElement('li');
        
        // Adaptable por si tu API usa minúsculas o mayúsculas en las propiedades del JSON
        const local = partido.equipo_local || partido.Local || partido.local;
        const visitante = partido.equipo_visitante || partido.Visitante || partido.visitante;
        const id = partido.id || partido.ID;

        li.innerHTML = `
            <span>${local} vs ${visitante}</span>
            <button class="btn-delete" data-id="${id}">Eliminar</button>
        `;
        
        listaPartidos.appendChild(li);
    });
}

// 2. Registrar un nuevo partido mediante el formulario (POST)
formPartido.addEventListener('submit', async (e) => {
    e.preventDefault(); // Evita que la página se recargue por defecto

    const nuevoPartido = {
        equipo_local: inputLocal.value.trim(),
        equipo_visitante: inputVisitante.value.trim()
    };

    try {
        const response = await fetch(API_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(nuevoPartido),
        });

        if (!response.ok) throw new Error('Error al registrar el partido');

        // Limpiar los inputs del formulario
        inputLocal.value = '';
        inputVisitante.value = '';

        // Refrescar la lista de partidos automáticamente
        obtenerPartidos();
    } catch (error) {
        console.error('Error en la petición POST:', error);
    }
});

// 3. Implementación de Delegación de Eventos para el Borrado (DELETE)
listaPartidos.addEventListener('click', async (e) => {
    if (e.target.classList.contains('btn-delete')) {
        const id = e.target.getAttribute('data-id');

        try {
            const response = await fetch(`${API_URL}/${id}`, {
                method: 'DELETE',
            });

            if (!response.ok) throw new Error('Error al eliminar el partido');

            // Refrescar la lista tras un borrado exitoso
            obtenerPartidos();
        } catch (error) {
            console.error('Error en la petición DELETE:', error);
        }
    }
});