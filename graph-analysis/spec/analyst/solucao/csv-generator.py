import csv
import os
import random
import schedule
import time
from datetime import datetime, timedelta

# Configurações de ambiente
NUM_COUNTRIES = int(os.getenv("GENERATOR_COUNTRY_COUNT", 30))
OUTPUT_DIR = os.getenv("OUTPUT_DIR", "/data")
INTERVAL_MINUTES = int(os.getenv("GENERATOR_INTERVAL_MINUTES", 5))
MAX_VACCINES_PER_COUNTRY = 3
START_DATE = datetime(2021, 1, 1)

READY_FLAG_PATH = os.path.join(OUTPUT_DIR, "ready.flag")

# Países ampliados
COUNTRIES = [
    ("BRA", "Brazil"), ("USA", "United States"), ("DEU", "Germany"),
    ("FRA", "France"), ("ITA", "Italy"), ("ESP", "Spain"),
    ("CAN", "Canada"), ("MEX", "Mexico"), ("ARG", "Argentina"),
    ("CHN", "China"), ("JPN", "Japan"), ("KOR", "South Korea"),
    ("IND", "India"), ("RUS", "Russia"), ("AUS", "Australia"),
    ("ZAF", "South Africa"), ("GBR", "United Kingdom"), ("IRL", "Ireland"),
    ("NLD", "Netherlands"), ("SWE", "Sweden"), ("NOR", "Norway"),
    ("DNK", "Denmark"), ("FIN", "Finland"), ("POL", "Poland"),
    ("TUR", "Turkey"), ("EGY", "Egypt"), ("SAU", "Saudi Arabia"),
    ("NZL", "New Zealand"), ("COL", "Colombia"), ("PER", "Peru")
]

VACCINES = [
    "Pfizer-BioNTech", "Moderna", "AstraZeneca", "Sinovac", "Sputnik V"
]

def safe_write_csv(tmp_path, final_path, write_fn):
    with open(tmp_path, "w", newline="") as f:
        writer = csv.writer(f)
        write_fn(writer)
    os.rename(tmp_path, final_path)

def generate_csvs():
    os.makedirs(OUTPUT_DIR, exist_ok=True)

    # Remove ready.flag
    if os.path.exists(READY_FLAG_PATH):
        os.remove(READY_FLAG_PATH)

    total_cases = 0
    total_vaccinations = 0
    total_country_vaccines = 0

    # countries.csv
    safe_write_csv(
        f"{OUTPUT_DIR}/countries.csv.tmp",
        f"{OUTPUT_DIR}/countries.csv",
        lambda writer: (
            writer.writerow(["iso3", "name"]),
            [writer.writerow([iso3, name]) for iso3, name in COUNTRIES[:NUM_COUNTRIES]]
        )
    )

    # covid_cases.csv
    def write_covid_cases(writer):
        nonlocal total_cases
        writer.writerow(["iso3", "date", "total_cases", "total_deaths"])
        for iso3, _ in COUNTRIES[:NUM_COUNTRIES]:
            cases_for_country = random.randint(8000, 12000)
            for i in range(cases_for_country):
                date = (START_DATE + timedelta(days=i)).strftime("%Y-%m-%d")
                cases = random.randint(10000, 500000)
                deaths = int(cases * random.uniform(0.01, 0.03))
                writer.writerow([iso3, date, cases, deaths])
                total_cases += 1

    safe_write_csv(
        f"{OUTPUT_DIR}/covid_cases.csv.tmp",
        f"{OUTPUT_DIR}/covid_cases.csv",
        write_covid_cases
    )

    # vaccinations.csv
    def write_vaccinations(writer):
        nonlocal total_vaccinations
        writer.writerow(["iso3", "date", "total_vaccinated"])
        for iso3, _ in COUNTRIES[:NUM_COUNTRIES]:
            vaccinations_for_country = random.randint(8000, 12000)
            for i in range(vaccinations_for_country):
                date = (START_DATE + timedelta(days=i)).strftime("%Y-%m-%d")
                vaccinated = random.randint(10000, 1000000)
                writer.writerow([iso3, date, vaccinated])
                total_vaccinations += 1

    safe_write_csv(
        f"{OUTPUT_DIR}/vaccinations.csv.tmp",
        f"{OUTPUT_DIR}/vaccinations.csv",
        write_vaccinations
    )

    # vaccines.csv
    safe_write_csv(
        f"{OUTPUT_DIR}/vaccines.csv.tmp",
        f"{OUTPUT_DIR}/vaccines.csv",
        lambda writer: (
            writer.writerow(["vaccine_name"]),
            [writer.writerow([name]) for name in VACCINES]
        )
    )

    # country_vaccines.csv
    def write_country_vaccines(writer):
        nonlocal total_country_vaccines
        writer.writerow(["iso3", "vaccine_name"])
        for iso3, _ in COUNTRIES[:NUM_COUNTRIES]:
            used = random.sample(VACCINES, random.randint(1, MAX_VACCINES_PER_COUNTRY))
            for vaccine in used:
                writer.writerow([iso3, vaccine])
                total_country_vaccines += 1

    safe_write_csv(
        f"{OUTPUT_DIR}/country_vaccines.csv.tmp",
        f"{OUTPUT_DIR}/country_vaccines.csv",
        write_country_vaccines
    )

    # vaccine_approvals.csv
    safe_write_csv(
        f"{OUTPUT_DIR}/vaccine_approvals.csv.tmp",
        f"{OUTPUT_DIR}/vaccine_approvals.csv",
        lambda writer: (
            writer.writerow(["vaccine_name", "date"]),
            [writer.writerow([
                name,
                (START_DATE - timedelta(days=random.randint(10, 100))).strftime("%Y-%m-%d")
            ]) for name in VACCINES]
        )
    )

    # Cria ready.flag
    with open(READY_FLAG_PATH, "w") as f:
        f.write("ready\n")

    now = datetime.now().isoformat(timespec="seconds")
    print(f"[{now}] CSVs gerados em '{OUTPUT_DIR}'")
    print(f"[{now}] Países gerados: {NUM_COUNTRIES}")
    print(f"[{now}] Casos de Covid gerados: {total_cases}")
    print(f"[{now}] Registros de vacinação gerados: {total_vaccinations}")
    print(f"[{now}] Relações país-vacina geradas: {total_country_vaccines}")
    print(f"[{now}] Vacinas registradas: {len(VACCINES)}")
    print(f"[{now}] Arquivo 'ready.flag' criado com sucesso.")

# Execução inicial e agendamento
generate_csvs()
schedule.every(INTERVAL_MINUTES).minutes.do(generate_csvs)

while True:
    schedule.run_pending()
    time.sleep(1)
