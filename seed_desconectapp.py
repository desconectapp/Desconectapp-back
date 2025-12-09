import requests
import random
from typing import Dict, List, Tuple

BASE_URL = "http://localhost:8080"

ADMIN_EMAIL = "admin@desconect.app"
ADMIN_PASSWORD = "password123"

UNSPLASH_ACCESS_KEY = "xIyhFleLUlZ-_2d2aovnO-nYP_zACF1dRebGif340KU"
UNSPLASH_BASE_URL = "https://api.unsplash.com"

# Cantidad de usuarios
TOTAL_USERS = 48

# Actividades "core" donde queremos MUCHOS matchings
MATCHING_ACTIVITIES = ["Swimming", "Biking", "Gym"]

# TODAS las actividades que vamos a crear (name, category)
TARGET_ACTIVITIES = [
    # Core deportivas
    ("Swimming", "SPORT"),
    ("Biking", "SPORT"),
    ("Gym", "SPORT"),

    # Más deporte
    ("Running", "SPORT"),
    ("Padel", "SPORT"),
    ("Tennis", "SPORT"),
    ("Climbing", "SPORT"),
    ("Crossfit", "SPORT"),
    ("Functional Training", "SPORT"),
    ("Boxing", "SPORT"),

    # Outdoor
    ("Hiking", "OUTDOOR"),
    ("Trail Running", "OUTDOOR"),
    ("Kayak", "OUTDOOR"),
    ("Outdoor Workout", "OUTDOOR"),

    # Wellness / bien estar
    ("Yoga", "WELLNESS"),
    ("Meditation", "WELLNESS"),
    ("Pilates", "WELLNESS"),
    ("Stretching", "WELLNESS"),

    # Social / juegos
    ("Board Games", "GAME"),
    ("Chess", "GAME"),
    ("Dungeons & Dragons", "GAME"),
    ("Fútbol 5", "SPORT"),
    ("Basket", "SPORT"),

    # Creativo / indoor
    ("Photography Walk", "CREATIVE"),
    ("Street Photography", "CREATIVE"),
    ("Painting Workshop", "CREATIVE"),
    ("Cooking Class", "INDOOR"),
    ("Baking & Coffee", "INDOOR"),

    # Social / idiomas / baile
    ("Language Exchange", "SOCIAL"),
    ("After Office", "SOCIAL"),
    ("Salsa", "SOCIAL"),
    ("Bachata", "SOCIAL"),
    ("K-Pop Dance", "SOCIAL"),
]

# Barrios porteños + coords aproximadas
BA_NEIGHBOURHOODS = [
    ("Palermo",   -34.57, -58.43),
    ("Belgrano",  -34.56, -58.45),
    ("Caballito", -34.61, -58.44),
    ("Almagro",   -34.61, -58.42),
    ("Recoleta",  -34.59, -58.39),
    ("San Telmo", -34.62, -58.37),
    ("La Boca",   -34.63, -58.36),
    ("Nuñez",     -34.54, -58.46),
    ("Villa Urquiza", -34.58, -58.49),
    ("Chacarita", -34.59, -58.45),
]

# Nombres argentinos
FIRST_NAMES = [
    "Lautaro", "Agustina", "Sofía", "Bruno", "Camila", "Juan", "Lucía",
    "Franco", "Micaela", "Nicolás", "Valentina", "Matías", "Florencia",
    "Diego", "Rocío", "Gonzalo", "Julieta", "Facundo", "Milagros",
    "Santiago", "Martina", "Ezequiel", "Carla", "Julián", "Malena",
    "Tomás", "Jimena", "Pedro", "Paula", "Ramiro", "Ayelén",
]

LAST_NAMES = [
    "González", "Rodríguez", "Fernández", "López", "Martínez", "Gómez",
    "Díaz", "Pérez", "Sánchez", "Romero", "Sosa", "Torres", "Ramírez",
    "Flores", "Acosta", "Rivas", "Benítez", "Herrera", "Molina", "Castro",
]

# Nombres raros para grupos
WEIRD_GROUP_NAMES = [
    "Cardio Apocalíptico",
    "Los Nadadores del Fin del Mundo",
    "Team Oxígeno Cero",
    "Pedaleando en Chanclas",
    "La Logia del Press Banca",
    "Piscina Intergaláctica",
    "Bicis & Fainá",
    "Crossfit Emocional",
    "Los de la Barra de Proteína",
    "Club de los Cuadríceps Tristes",
]

# Timeslot clusters para generar muchos matchings
TIMESLOT_CLUSTERS = {
    "swimming_evening": [100, 101, 102],
    "swimming_morning": [50, 51, 52],
    "biking_morning":   [60, 61, 62],
    "biking_afternoon": [180, 181, 182],
    "gym_after_office": [220, 221, 222],
    "gym_night":        [280, 281, 282],
}

# Clusters genéricos por categoría para las actividades "extra"
CATEGORY_DEFAULT_CLUSTERS: Dict[str, List[List[int]]] = {
    "SPORT":    [[300, 301], [302, 303]],
    "OUTDOOR":  [[310, 311], [312, 313]],
    "INDOOR":   [[320, 321]],
    "GAME":     [[330, 331]],
    "CREATIVE": [[340, 341]],
    "SOCIAL":   [[350, 351], [352, 353]],
    "WELLNESS": [[360, 361], [362, 363]],
}


# ==========================
# HTTP helpers
# ==========================

def api_post(path: str, token: str = None, json_data=None):
    headers = {}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    resp = requests.post(f"{BASE_URL}{path}", json=json_data, headers=headers)
    if resp.status_code >= 400:
        raise RuntimeError(f"POST {path} failed [{resp.status_code}]: {resp.status_code} {resp.text}")
    return resp.json() if resp.text else None


def api_get(path: str, token: str = None, params=None):
    headers = {}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    resp = requests.get(f"{BASE_URL}{path}", headers=headers, params=params)
    if resp.status_code >= 400:
        raise RuntimeError(f"GET {path} failed [{resp.status_code}]: {resp.status_code} {resp.text}")
    return resp.json() if resp.text else None


def api_put(path: str, token: str = None, json_data=None):
    headers = {}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    resp = requests.put(f"{BASE_URL}{path}", json=json_data, headers=headers)
    if resp.status_code >= 400:
        raise RuntimeError(f"PUT {path} failed [{resp.status_code}]: {resp.status_code} {resp.text}")
    return resp.json() if resp.text else None


# ==========================
# Auth & Activities
# ==========================

def login_admin() -> str:
    print(f"[+] Logueando admin {ADMIN_EMAIL}")
    data = {"email": ADMIN_EMAIL, "password": ADMIN_PASSWORD}
    res = api_post("/auth/login", json_data=data)
    token = res["token"]
    print("[+] Admin logueado OK")
    return token


def get_or_create_activity_id(token: str, name: str, category: str) -> int:
    print(f"[+] Buscando actividad '{name}'")
    activities = api_get(
        "/admin/activities",
        token=token,
        params={"q": name, "_start": 0, "_end": 100},
    )
    if isinstance(activities, list) and activities:
        for act in activities:
            if act.get("name") == name:
                act_id = act["id"]
                print(f"    -> Encontrada existente: id={act_id}")
                return act_id

    print(f"    -> No existe, creando '{name}' en categoría {category}")
    payload = {
        "name": name,
        "icon": None,
        "category": category,
    }
    created = api_post("/admin/activities", token=token, json_data=payload)
    act_id = created["id"]
    print(f"    -> Creada actividad id={act_id}")
    return act_id


# ==========================
# Users
# ==========================

def random_argentinian_user(i: int) -> Tuple[str, str, str]:
    first = random.choice(FIRST_NAMES)
    last = random.choice(LAST_NAMES)
    full_name = f"{first} {last}"
    email_slug = (
        full_name.lower()
        .replace(" ", ".")
        .replace("á", "a")
        .replace("é", "e")
        .replace("í", "i")
        .replace("ó", "o")
        .replace("ú", "u")
        .replace("ñ", "n")
    )
    email = f"{email_slug}{i}@desconectapp.test"
    password = "password123"
    return full_name, email, password


def create_users_via_admin(token: str, total: int) -> List[int]:
    print(f"[+] Creando {total} usuarios vía /admin/users")
    user_ids = []
    for i in range(total):
        name, email, password = random_argentinian_user(i)
        payload = {
            "name": name,
            "email": email,
            "password": password,
            "is_admin": False,
        }
        res = api_post("/admin/users", token=token, json_data=payload)
        user_id = res["id"]
        user_ids.append(user_id)
        print(f"    -> Usuario {i+1}/{total}: {name} (id={user_id}, email={email})")
    return user_ids


# ==========================
# Activity Requests
# ==========================

def pick_neighbourhood() -> Tuple[str, float, float]:
    name, lat, lng = random.choice(BA_NEIGHBOURHOODS)
    lat += random.uniform(-0.005, 0.005)
    lng += random.uniform(-0.005, 0.005)
    return name, lat, lng


def create_activity_request_for_user(
    admin_token: str,
    user_id: int,
    activity_id: int,
    description: str,
    timeslots: List[int],
):
    location_name, lat, lng = pick_neighbourhood()

    payload = {
        "user_id": user_id,
        "activity_id": activity_id,
        "description": description,
        "longitude": lng,
        "latitude": lat,
        "search_radius": 8,
        "max_participants": 6,
        "participants_needed": 3,
        "timeslots": timeslots,
    }

    res = api_post("/activities/request", token=admin_token, json_data=payload)
    req_id = res["id"]
    print(
        f"    -> ActivityRequest id={req_id} user_id={user_id} "
        f"activity_id={activity_id} barrio={location_name}"
    )


def seed_core_match_activities(
    admin_token: str,
    activity_ids: Dict[str, int],
    user_ids: List[int],
):
    """
    Genera muchos matchings solo para Swimming / Biking / Gym.
    """
    print("[+] Creando activity_requests fuertes para Swimming / Biking / Gym")

    random.shuffle(user_ids)
    chunk_size = len(user_ids) // 3
    users_swimming = user_ids[0:chunk_size]
    users_biking = user_ids[chunk_size: 2 * chunk_size]
    users_gym = user_ids[2 * chunk_size:]

    swim_id = activity_ids["Swimming"]
    bike_id = activity_ids["Biking"]
    gym_id = activity_ids["Gym"]

    # Swimming
    swimming_clusters = [
        TIMESLOT_CLUSTERS["swimming_morning"],
        TIMESLOT_CLUSTERS["swimming_evening"],
    ]
    print("  [*] Swimming")
    for idx, user_id in enumerate(users_swimming):
        cluster = swimming_clusters[idx % len(swimming_clusters)]
        desc = random.choice([
            "Nadar un rato y tomar mate después",
            "Entrenamiento suave de pileta",
            "Natación para despejar la cabeza",
            "Series de crawl y algo de pecho",
        ])
        create_activity_request_for_user(admin_token, user_id, swim_id, desc, cluster)

    # Biking
    biking_clusters = [
        TIMESLOT_CLUSTERS["biking_morning"],
        TIMESLOT_CLUSTERS["biking_afternoon"],
    ]
    print("  [*] Biking")
    for idx, user_id in enumerate(users_biking):
        cluster = biking_clusters[idx % len(biking_clusters)]
        desc = random.choice([
            "Salidita tranqui por los bosques de Palermo",
            "Vuelta larga por Costanera",
            "Entrenamiento de bici para fondo",
            "Rodaje suave con algo de desnivel",
        ])
        create_activity_request_for_user(admin_token, user_id, bike_id, desc, cluster)

    # Gym
    gym_clusters = [
        TIMESLOT_CLUSTERS["gym_after_office"],
        TIMESLOT_CLUSTERS["gym_night"],
    ]
    print("  [*] Gym")
    for idx, user_id in enumerate(users_gym):
        cluster = gym_clusters[idx % len(gym_clusters)]
        desc = random.choice([
            "Rutina de fuerza full-body",
            "Pierna y algo de cardio",
            "Pecho, espalda y un poco de bici fija",
            "Gimnasio post laburo para descargar",
        ])
        create_activity_request_for_user(admin_token, user_id, gym_id, desc, cluster)

    print("[+] Activity requests core creadas – el matcher debería armar varios grupos grandes")


def seed_extra_activities(
    admin_token: str,
    activity_ids: Dict[str, int],
    activity_categories: Dict[str, str],
    user_ids: List[int],
):
    """
    Para cada actividad que no sea Swimming/Biking/Gym,
    crea algunos activity_requests más livianos para poblar el sistema.
    """
    print("[+] Creando activity_requests para actividades extra")

    # No queremos abusar del mismo user; elegimos random 3–5 users por actividad
    for activity_name, act_id in activity_ids.items():
        if activity_name not in MATCHING_ACTIVITIES:
            category = activity_categories.get(activity_name, "SOCIAL")
            clusters = CATEGORY_DEFAULT_CLUSTERS.get(category, [[400, 401]])
            # Entre 3 y 5 usuarios por actividad
            num_users = random.randint(3, 5)
            users_for_activity = random.sample(user_ids, min(num_users, len(user_ids)))

            print(f"  [*] {activity_name} ({category}) - {len(users_for_activity)} usuarios")
            for idx, uid in enumerate(users_for_activity):
                cluster = clusters[idx % len(clusters)]
                desc = f"{activity_name} en grupo por CABA"
                create_activity_request_for_user(
                    admin_token,
                    uid,
                    act_id,
                    desc,
                    cluster,
                )


# ==========================
# Grupos: listar & nombre raro
# ==========================

def fetch_all_groups_admin(admin_token: str) -> List[Dict]:
    """
    Usa /admin/groups con un rango grande para traer todos los grupos.
    """
    print("[+] Obteniendo todos los grupos via /admin/groups")
    groups = api_get(
        "/admin/groups",
        token=admin_token,
        params={"_start": 0, "_end": 1000},
    )
    if not isinstance(groups, list):
        raise RuntimeError("Respuesta inesperada de /admin/groups (no es lista)")
    print(f"    -> {len(groups)} grupos encontrados")
    return groups


def rename_some_groups_with_weird_names(admin_token: str, groups: List[Dict]):
    print("[+] Renombrando algunos grupos con nombres raros")
    group_ids = [g.get("id") or g.get("ID") for g in groups if g.get("id") or g.get("ID")]
    random.shuffle(group_ids)

    num_to_rename = max(1, len(group_ids) // 3)
    target_ids = group_ids[:num_to_rename]

    for gid in target_ids:
        new_name = random.choice(WEIRD_GROUP_NAMES)
        print(f"    -> Renombrando grupo {gid} a '{new_name}'")
        payload = {"name": new_name}
        api_put(f"/admin/groups/{gid}", token=admin_token, json_data=payload)

    print(f"[+] Renombrados {len(target_ids)} grupos con nombres raros")


# ==========================
# Unsplash + actualizar grupos (imagen + public)
# ==========================

def fetch_unsplash_images(query: str, per_page: int = 30) -> List[str]:
    """
    Obtiene URLs de imágenes de Unsplash usando la API de search/photos.
    Usa orientation=squarish como pediste.
    """
    print(f"[+] Buscando imágenes en Unsplash: query='{query}'")
    params = {
        "query": query,
        "orientation": "squarish",
        "per_page": per_page,
        "client_id": UNSPLASH_ACCESS_KEY,
    }

    resp = requests.get(f"{UNSPLASH_BASE_URL}/search/photos", params=params)
    if resp.status_code >= 400:
        raise RuntimeError(f"Unsplash search failed [{resp.status_code}]: {resp.status_code} {resp.text}")
    data = resp.json()
    results = data.get("results", [])
    urls = [r["urls"]["thumb"] for r in results if "urls" in r and "thumb" in r["urls"]]
    print(f"    -> {len(urls)} imágenes obtenidas")
    return urls




def update_groups_with_images_and_public(
    admin_token: str,
    groups: List[Dict],
):
    """
    Usa PUT /admin/groups/:id para:
      - setear avatar_url con una imagen de Unsplash
      - setear public = true
      - asegurar que TODO grupo tenga nombre:
          * si ya tiene name -> se respeta
          * si no tiene:
              - ~1/3 reciben un nombre raro de WEIRD_GROUP_NAMES
              - el resto "Grupo #<id>"
      - usar una búsqueda distinta en Unsplash por cada activity_id
        para que la foto represente la actividad del grupo.
    """

    # 1) Mapear activity_id -> nombre de la actividad
    print("[+] Construyendo mapa activity_id -> activity_name desde /admin/activities")
    activities = api_get(
        "/admin/activities",
        token=admin_token,
        params={"_start": 0, "_end": 1000},
    )
    activity_name_by_id: Dict[int, str] = {}
    if isinstance(activities, list):
        for a in activities:
            aid = a.get("id") or a.get("ID")
            aname = (a.get("name") or a.get("Name") or "").strip()
            if aid is not None and aname:
                activity_name_by_id[int(aid)] = aname

    # 2) Juntar todos los activity_id presentes en los grupos
    activity_ids_in_groups: set[int] = set()
    for g in groups:
        # Ajustá estas keys si tu JSON cambia
        aid = g.get("activity_id") or g.get("ActivityId") or g.get("activityId")
        if aid is not None:
            try:
                activity_ids_in_groups.add(int(aid))
            except (TypeError, ValueError):
                pass

    # 3) Para cada activity_id, pedir imágenes a Unsplash usando el nombre como query
    print("[+] Buscando imágenes en Unsplash por actividad")
    images_by_activity: Dict[int, List[str]] = {}
    for aid in activity_ids_in_groups:
        activity_name = activity_name_by_id.get(aid)
        if not activity_name:
            # Fallback genérico si no conocemos el nombre
            query = "fitness"
        else:
            query = activity_name

        print(f"    -> activity_id={aid}, query='{query}'")
        try:
            urls = fetch_unsplash_images(query, per_page=20)
        except Exception as e:
            print(f"       [WARN] Error buscando imágenes para '{query}': {e}")
            urls = []

        # Si no hay resultados, fallback a algo general
        if not urls:
            try:
                urls = fetch_unsplash_images("gym", per_page=20)
            except Exception as e:
                print(f"       [WARN] Fallback 'gym' también falló: {e}")
                urls = []

        if urls:
            images_by_activity[aid] = urls
        else:
            # último fallback: no tenemos imágenes, lo manejamos más abajo
            print(f"       [WARN] Sin imágenes disponibles para activity_id={aid}")

    if not images_by_activity:
        print("[WARN] No se pudieron obtener imágenes de Unsplash para ninguna actividad; salteando update de imágenes.")
        return

    print("[+] Actualizando grupos con avatar_url + public=true y nombres via /admin/groups/:id")

    # 4) Preparar info de nombres (quién ya tiene / quién no)
    group_info = []
    for g in groups:
        gid = g.get("id") or g.get("ID")
        if gid is None:
            continue
        current_name = (g.get("name") or g.get("Name") or "").strip()
        group_info.append((gid, current_name))

    unnamed_ids = [gid for (gid, current_name) in group_info if not current_name]

    weird_target_ids: set[int] = set()
    if unnamed_ids:
        random.shuffle(unnamed_ids)
        weird_count = max(1, len(unnamed_ids) // 3)
        weird_count = min(weird_count, len(unnamed_ids))
        weird_target_ids = set(unnamed_ids[:weird_count])

    # 5) Llevamos un índice de rotación de imágenes por activity_id
    img_idx_by_activity: Dict[int, int] = {}

    for g in groups:
        gid = g.get("id") or g.get("ID")
        if gid is None:
            continue

        # leemos activity_id del grupo
        aid_raw = g.get("activity_id") or g.get("ActivityId") or g.get("activityId")
        try:
            aid = int(aid_raw) if aid_raw is not None else None
        except (TypeError, ValueError):
            aid = None

        # elegimos la lista de imágenes correspondiente a esa actividad
        urls_for_activity: List[str] = []
        if aid is not None and aid in images_by_activity:
            urls_for_activity = images_by_activity[aid]

        # si por alguna razón no tenemos imágenes para esa actividad, hacemos fallback a cualquier otra
        if not urls_for_activity:
            any_aid = next(iter(images_by_activity.keys()))
            urls_for_activity = images_by_activity[any_aid]
            aid_for_index = any_aid
        else:
            aid_for_index = aid

        # índice de rotación por actividad
        current_idx = img_idx_by_activity.get(aid_for_index, 0)
        img_url = urls_for_activity[current_idx % len(urls_for_activity)]
        img_idx_by_activity[aid_for_index] = current_idx + 1

        current_name = (g.get("name") or g.get("Name") or "").strip()
        payload = {
            "avatar_url": img_url,
            "public": True,
        }

        # Sólo tocamos el nombre si estaba vacío
        if not current_name:
            if gid in weird_target_ids and WEIRD_GROUP_NAMES:
                new_name = random.choice(WEIRD_GROUP_NAMES)
            else:
                new_name = f"Grupo #{gid}"
            payload["name"] = new_name
            print(
                f"    -> Grupo {gid}: activity_id={aid} avatar_url={img_url[:60]}..., "
                f"public=True, name='{new_name}' (antes sin nombre)"
            )
        else:
            print(
                f"    -> Grupo {gid}: activity_id={aid} avatar_url={img_url[:60]}..., "
                f"public=True, name existente='{current_name}'"
            )

        api_put(f"/admin/groups/{gid}", token=admin_token, json_data=payload)

    print("[+] Todos los grupos actualizados: imagen por actividad, públicos y con nombre garantizado")





# ==========================
# Main
# ==========================

def main():
    admin_token = login_admin()
    # 1) Crear TODAS las actividades
    activity_ids: Dict[str, int] = {}
    activity_categories: Dict[str, str] = {}
    for name, category in TARGET_ACTIVITIES:
        act_id = get_or_create_activity_id(admin_token, name, category)
        activity_ids[name] = act_id
        activity_categories[name] = category

    # 2) Crear usuarios argentinos
    user_ids = create_users_via_admin(admin_token, TOTAL_USERS)

    # 3) Match fuerte para Swimming / Biking / Gym
    seed_core_match_activities(admin_token, activity_ids, user_ids)

    # 4) Activity requests livianos para TODAS las otras actividades
    seed_extra_activities(admin_token, activity_ids, activity_categories, user_ids)

    # 5) Obtener grupos actuales (los que generó tu matcher)
    groups = fetch_all_groups_admin(admin_token)

    # 6) Poner imagen + public=true y de paso renombrar algunos grupos
    update_groups_with_images_and_public(admin_token, groups)

    print("\n[✔] Seeding completo.")
    print("    - Usuarios creados:", len(user_ids))
    print("    - Actividades creadas:", len(activity_ids))
    print("    - Grupos totales:", len(groups))


if __name__ == "__main__":
    main()
