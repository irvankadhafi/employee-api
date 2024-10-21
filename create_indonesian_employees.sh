#!/bin/bash

# API endpoint
API_URL="http://localhost:8080/api/employees"

# Arrays of Indonesian names and tech-related positions
names=("Budi Santoso" "Agus Wijaya" "Siti Nurhaliza" "Ahmad Rian" "Dewi Lestari" "Rini Handayani" "Bagus Pratama" "Dian Sastro" "Taufik Hidayat" "Indra Gunawan" "Lestari Yulianti" "Yusuf Maulana" "Eka Putra" "Sri Lestari" "Nur Aisyah" "Rudy Hartono" "Febrianto Surya" "Cindy Rahma" "Wulan Pertiwi" "Agung Nugroho")
positions=("Software Engineer" "Backend Engineer" "Frontend Engineer" "Fullstack Developer" "DevOps Engineer" "Product Manager" "Data Analyst" "UX Designer" "QA Engineer" "CTO" "Mobile Developer" "Data Scientist" "Project Manager" "UI Designer" "Security Engineer" "Cloud Engineer" "Tech Lead" "Scrum Master" "Database Administrator" "Systems Architect")

# Function to generate random salary between 5,000,000 and 20,000,000 IDR
generate_salary() {
  echo $((5000000 + RANDOM % 15000001))
}

# Loop to create 100 employees
for i in {1..100}; do
  # Select random name and position
  name=${names[$RANDOM % ${#names[@]}]}
  position=${positions[$RANDOM % ${#positions[@]}]}
  salary=$(generate_salary)

  # Send POST request to create the employee
  curl -X POST $API_URL \
    -H "Content-Type: application/json" \
    -d '{
      "name": "'"$name"'",
      "position": "'"$position"'",
      "salary": '"$salary"'
    }'

  echo "Created employee $i: $name - $position with salary $salary IDR"
done