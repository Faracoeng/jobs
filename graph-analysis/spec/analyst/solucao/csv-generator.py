import csv
import os
import random
import schedule
import time
from datetime import datetime, timedelta

# Configuráveis por variável de ambiente
NUM_COUNTRIES = int(os.getenv("GENERATOR_COUNTRY_COUNT", 3))
NUM_CASES_PER_COUNTRY = int(os.getenv("GENERATOR_CASES_PER_COUNTRY", 10))
OUTPUT_DIR = os.getenv("OUTPUT_DIR", "data")
INTERVAL_MINUTES = int(os.getenv("GENERATOR_INTERVAL_MINUTES", 10))

COUNTRIES = [
    ("BRA", "Brazil"),
    ("USA", "United States"),
    ("DEU", "Germany"),
    ("FRA", "France"),
    ("ITA", "Italy"),
    ("ESP", "Spain")
]

VACCINES = [
    "Pfizer-BioNTech",
    "Moderna",
    "AstraZeneca",
    "Sinovac",
    "Sputnik V"
]

START_DATE = datetime(2021, 1, 1)

def generate_csvs():
    os.makedirs(OUTPUT_DIR, exist_ok=True)

    # country.csv
    with open(f"{OUTPUT_DIR}/countries.csv", "w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["iso3", "name"])
        for iso3, name in COUNTRIES[:NUM_COUNTRIES]:
            writer.writerow([iso3, name])

    # covid_cases.csv
    with open(f"{OUTPUT_DIR}/covid_cases.csv", "w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["iso3", "date", "total_cases", "total_deaths"])
        for iso3, _ in COUNTRIES[:NUM_COUNTRIES]:
            for i in range(NUM_CASES_PER_COUNTRY):
                date = (START_DATE + timedelta(days=i * 7)).strftime("%Y-%m-%d")
                cases = random.randint(10000, 1000000)
                deaths = int(cases * 0.02)
                writer.writerow([iso3, date, cases, deaths])

    # vaccinations.csv
    with open(f"{OUTPUT_DIR}/vaccinations.csv", "w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["iso3", "date", "total_vaccinated"])
        for iso3, _ in COUNTRIES[:NUM_COUNTRIES]:
            for i in range(NUM_CASES_PER_COUNTRY):
                date = (START_DATE + timedelta(days=i * 7)).strftime("%Y-%m-%d")
                vaccinated = random.randint(10000, 1000000)
                writer.writerow([iso3, date, vaccinated])

    # vaccines.csv
    with open(f"{OUTPUT_DIR}/vaccines.csv", "w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["vaccine_name"])
        for name in VACCINES:
            writer.writerow([name])

    # country_vaccines.csv
    with open(f"{OUTPUT_DIR}/country_vaccines.csv", "w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["iso3", "vaccine_name"])
        for iso3, _ in COUNTRIES[:NUM_COUNTRIES]:
            used = random.sample(VACCINES, random.randint(1, 3))
            for vaccine in used:
                writer.writerow([iso3, vaccine])

    # vaccine_approvals.csv
    with open(f"{OUTPUT_DIR}/vaccine_approvals.csv", "w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["vaccine_name", "date"])
        for name in VACCINES:
            date = (START_DATE - timedelta(days=random.randint(1, 60))).strftime("%Y-%m-%d")
            writer.writerow([name, date])

    print(f"[{datetime.now().isoformat()}] CSVs gerados em '{OUTPUT_DIR}'")

# Execução inicial e agendamento
generate_csvs()
schedule.every(INTERVAL_MINUTES).minutes.do(generate_csvs)

while True:
    schedule.run_pending()
    time.sleep(1)
