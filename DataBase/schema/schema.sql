--1. Ligas
CREATE TABLE ligas(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    pais VARCHAR(50) NOT NULL
);

--2.EQUIPOS
CREATE TABLE equipos(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE
);

--3. TABLAS DE POSICIONES
CREATE TABLE tabla_posiciones(
    id SERIAL PRIMARY KEY,
    liga_id INT NOT NULL REFERENCES ligas(id) ON DELETE CASCADE,
    equipo_id INT NOT NULL REFERENCES equipos(id) ON DELETE CASCADE,
    puntos INT DEFAULT 0,
    partidos_jugados INT DEFAULT 0,
    partidos_ganados INT DEFAULT 0,
    partidos_empatados INT DEFAULT 0,
    partidos_perdidos INT DEFAULT 0,

    UNIQUE (liga_id, equipo_id)
);

--4. PARTIDOS
CREATE TABLE partidos(
    id SERIAL PRIMARY KEY,
    liga_id INT NOT NULL REFERENCES ligas(id) ON DELETE CASCADE,
    equipo_local_id INT NOT NULL REFERENCES equipos(id),
    equipo_visitante_id INT NOT NULL REFERENCES equipos(id),
    goles_local INT,
    goles_visitante INT,
    estado VARCHAR(20) DEFAULT 'pendiente'
);
