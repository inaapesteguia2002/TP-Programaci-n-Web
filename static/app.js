const URL_LIGAS = 'http://localhost:8080/ligas';
const URL_EQUIPOS = 'http://localhost:8080/equipos';
const URL_PARTIDOS = 'http://localhost:8080/partidos';
const URL_TABLA = 'http://localhost:8080/tabla';

// Carga inicial de datos
document.addEventListener('DOMContentLoaded', () => {
    obtenerLigas();
    obtenerEquipos();
    obtenerPartidos();
});

// 1. LIGAS
async function obtenerLigas() {
    try {
        const res = await fetch(URL_LIGAS);
        const ligas = await res.json();
        
        const listaLigas = document.getElementById('lista-ligas');
        const selectLigaPartido = document.getElementById('select-liga'); 
        const selectLigaTabla = document.getElementById('select-liga-tabla'); 
        const selectVerTabla = document.getElementById('select-ver-tabla'); 
        
        listaLigas.innerHTML = '';
        selectLigaPartido.innerHTML = '<option value="">Seleccione una Liga...</option>';
        selectLigaTabla.innerHTML = '<option value="">Seleccione Liga...</option>';
        selectVerTabla.innerHTML = '<option value="">Ver posiciones de la Liga...</option>';
        
        if (!ligas) return;

        ligas.forEach(liga => {
            const id = liga.id || liga.ID;
            const nombre = liga.nombre || liga.Nombre;
            
            const li = document.createElement('li');
            li.innerHTML = `<span><strong>${nombre}</strong> (${liga.pais || liga.Pais})</span> <button class="btn-delete" data-id="${id}" data-tipo="liga">Eliminar</button>`;
            listaLigas.appendChild(li);

            const optionHTML = `<option value="${id}">${nombre}</option>`;
            selectLigaPartido.insertAdjacentHTML('beforeend', optionHTML);
            selectLigaTabla.insertAdjacentHTML('beforeend', optionHTML);
            selectVerTabla.insertAdjacentHTML('beforeend', optionHTML);
        });
    } catch (error) { console.error('Error Ligas:', error); }
}

document.getElementById('form-liga').addEventListener('submit', async (e) => {
    e.preventDefault();
    const body = {
        nombre: document.getElementById('nombre-liga').value.trim(),
        pais: document.getElementById('pais-liga').value.trim()
    };
    await fetch(URL_LIGAS, { method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(body) });
    document.getElementById('form-liga').reset();
    obtenerLigas();
});

// 2. EQUIPOS
async function obtenerEquipos() {
    try {
        const res = await fetch(URL_EQUIPOS);
        const equipos = await res.json();
        
        const listaEquipos = document.getElementById('lista-equipos');
        const selectLocal = document.getElementById('select-local');
        const selectVisitante = document.getElementById('select-visitante');
        const selectEquipoTabla = document.getElementById('select-equipo-tabla');
        
        listaEquipos.innerHTML = '';
        selectLocal.innerHTML = '<option value="">Equipo Local...</option>';
        selectVisitante.innerHTML = '<option value="">Equipo Visitante...</option>';
        selectEquipoTabla.innerHTML = '<option value="">Seleccione Equipo...</option>';
        
        if (!equipos) return;

        equipos.forEach(eq => {
            const id = eq.id || eq.ID;
            const nombre = eq.nombre || eq.Nombre;
            
            const li = document.createElement('li');
            li.innerHTML = `<span>${nombre}</span> <button class="btn-delete" data-id="${id}" data-tipo="equipo">Eliminar</button>`;
            listaEquipos.appendChild(li);

            const optionHTML = `<option value="${id}">${nombre}</option>`;
            selectLocal.insertAdjacentHTML('beforeend', optionHTML);
            selectVisitante.insertAdjacentHTML('beforeend', optionHTML);
            selectEquipoTabla.insertAdjacentHTML('beforeend', optionHTML);
        });
    } catch (error) { console.error('Error Equipos:', error); }
}

document.getElementById('form-equipo').addEventListener('submit', async (e) => {
    e.preventDefault();
    const body = { nombre: document.getElementById('nombre-equipo').value.trim() };
    await fetch(URL_EQUIPOS, { method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(body) });
    document.getElementById('form-equipo').reset();
    obtenerEquipos();
});

// Helper para extraer valores que pueden venir en mayúscula o como sql.NullInt32 / sql.NullString
function extraerValor(campo, valorPorDefecto = 0) {
    if (campo === undefined || campo === null) return valorPorDefecto;
    if (typeof campo === 'object') {
        if ('Int32' in campo) return campo.Valid ? campo.Int32 : valorPorDefecto;
        if ('String' in campo) return campo.Valid ? campo.String : valorPorDefecto;
    }
    return campo;
}

// 3. PARTIDOS
async function obtenerPartidos() {
    try {
        const res = await fetch(URL_PARTIDOS);
        const partidos = await res.json();
        const listaPartidos = document.getElementById('lista-partidos');
        listaPartidos.innerHTML = '';
        
        if (!partidos) return;

        partidos.forEach(p => {
            const id = p.id || p.ID;
            const localId = p.equipo_local_id || p.EquipoLocalID;
            const visitanteId = p.equipo_visitante_id || p.EquipoVisitanteID;
            const golesLocal = extraerValor(p.goles_local ?? p.GolesLocal, 0);
            const golesVisitante = extraerValor(p.goles_visitante ?? p.GolesVisitante, 0);
            const estado = extraerValor(p.estado ?? p.Estado, 'Pendiente');

            const li = document.createElement('li');
            li.innerHTML = `
                <span>Eq. ${localId} <strong>(${golesLocal} - ${golesVisitante})</strong> Eq. ${visitanteId} <br><small>${estado}</small></span>
                <button class="btn-delete" data-id="${id}" data-tipo="partido">Eliminar</button>
            `;
            listaPartidos.appendChild(li);
        });
    } catch (error) { console.error('Error Partidos:', error); }
}

document.getElementById('form-partido').addEventListener('submit', async (e) => {
    e.preventDefault();
    const localId = parseInt(document.getElementById('select-local').value);
    const visitanteId = parseInt(document.getElementById('select-visitante').value);

    if (localId === visitanteId) {
        alert("El equipo local no puede ser igual al visitante."); 
        return;
    }

    const nuevoPartido = {
        liga_id: parseInt(document.getElementById('select-liga').value),
        equipo_local_id: localId,
        equipo_visitante_id: visitanteId,
        goles_local: parseInt(document.getElementById('goles-local').value),
        goles_visitante: parseInt(document.getElementById('goles-visitante').value),
        estado: document.getElementById('estado-partido').value.trim()
    };

    await fetch(URL_PARTIDOS, { method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(nuevoPartido) });
    document.getElementById('form-partido').reset();
    obtenerPartidos();
});

// 4. TABLA DE POSICIONES
document.getElementById('form-tabla').addEventListener('submit', async (e) => {
    e.preventDefault();
    const inscripcion = {
        liga_id: parseInt(document.getElementById('select-liga-tabla').value),
        equipo_id: parseInt(document.getElementById('select-equipo-tabla').value)
    };

    try {
        const res = await fetch(URL_TABLA, { method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify(inscripcion) });
        if (!res.ok) {
            const errorText = await res.text();
            alert("Error al inscribir: " + errorText);
            return;
        }
        
        document.getElementById('form-tabla').reset();
        
        const ligaActual = document.getElementById('select-ver-tabla').value;
        if (ligaActual == inscripcion.liga_id) {
            verTabla(ligaActual);
        } else {
            alert("Equipo inscripto con éxito.");
        }
    } catch (error) { console.error('Error al inscribir equipo:', error); }
});

document.getElementById('select-ver-tabla').addEventListener('change', async (e) => {
    const ligaId = e.target.value;
    if (ligaId) verTabla(ligaId);
    else document.getElementById('lista-tabla').innerHTML = '<li style="text-align: center; color: #888; font-size: 13px;">Selecciona una liga para ver la tabla.</li>';
});

async function verTabla(ligaId) {
    try {
        const res = await fetch(`${URL_TABLA}/${ligaId}`);
        const tabla = await res.json();
        const listaTabla = document.getElementById('lista-tabla');
        listaTabla.innerHTML = '';

        if (!tabla || tabla.length === 0) {
            listaTabla.innerHTML = '<li>No hay equipos en esta liga.</li>';
            return;
        }

        tabla.forEach(pos => {
            const lId = pos.liga_id || pos.LigaID;
            const eqId = pos.equipo_id || pos.EquipoID;
            const pts = extraerValor(pos.puntos ?? pos.Puntos, 0);
            const pg = extraerValor(pos.partidos_ganados ?? pos.PartidosGanados, 0);
            const pe = extraerValor(pos.partidos_empatados ?? pos.PartidosEmpatados, 0);
            const pp = extraerValor(pos.partidos_perdidos ?? pos.PartidosPerdidos, 0);

            const li = document.createElement('li');
            li.innerHTML = `
                <span>Eq. ${eqId} | <strong>Pts: ${pts}</strong> | G:${pg} E:${pe} P:${pp}</span>
                <button class="btn-delete" data-ligaid="${lId}" data-equipoid="${eqId}" data-tipo="tabla">Sacar</button>
            `;
            listaTabla.appendChild(li);
        });
    } catch (error) { console.error('Error al traer la tabla:', error); }
}

// 5. DELEGACIÓN DE EVENTOS GLOBAL (DELETE)
document.addEventListener('click', async (e) => {
    if (e.target.classList.contains('btn-delete')) {
        const tipo = e.target.getAttribute('data-tipo');
        
        // El DELETE de la tabla es especial porque manda JSON en el body
        if (tipo === 'tabla') {
            const ligaId = parseInt(e.target.getAttribute('data-ligaid'));
            const equipoId = parseInt(e.target.getAttribute('data-equipoid'));
            
            await fetch(URL_TABLA, { 
                method: 'DELETE', 
                headers: {'Content-Type': 'application/json'}, 
                body: JSON.stringify({ liga_id: ligaId, equipo_id: equipoId }) 
            });
            verTabla(ligaId);
            return;
        }

        // DELETE normal por parámetro en la URL
        const id = e.target.getAttribute('data-id');
        let urlTarget = '';
        if (tipo === 'liga') urlTarget = `${URL_LIGAS}/${id}`;
        if (tipo === 'equipo') urlTarget = `${URL_EQUIPOS}/${id}`;
        if (tipo === 'partido') urlTarget = `${URL_PARTIDOS}/${id}`;

        try {
            await fetch(urlTarget, { method: 'DELETE' });
            if (tipo === 'liga') obtenerLigas();
            if (tipo === 'equipo') obtenerEquipos();
            if (tipo === 'partido') obtenerPartidos();
        } catch (error) {
            console.error('Error al eliminar:', error);
        }
    }
});
