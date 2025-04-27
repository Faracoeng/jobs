package entity


import "time"


type Country struct {
   ISO3 string
   Name string
}


type Vaccine struct {
   Name string
}


type CovidCase struct {
   ISO3        string
   Date        time.Time
   TotalCases  int
   TotalDeaths int
}


type VaccinationStat struct {
   ISO3            string
   Date            time.Time
   TotalVaccinated int
}


type VaccineApproval struct {
   Date        time.Time
}


type CountryVaccine struct {
   ISO3        string
   VaccineName string
}