// Node Keys
CREATE CONSTRAINT IF NOT EXISTS ON (n:Vaccine) ASSERT (n.product) IS NODE KEY;
CREATE CONSTRAINT IF NOT EXISTS ON (n:Date) ASSERT (n.date) IS NODE KEY;
CREATE CONSTRAINT IF NOT EXISTS ON (n:Region) ASSERT (n.name) IS NODE KEY;
CREATE CONSTRAINT IF NOT EXISTS ON (n:Country) ASSERT (n.code) IS NODE KEY;

// Node Properties Must Exist
CREATE CONSTRAINT IF NOT EXISTS ON (n:Vaccine) ASSERT exists(n.company);
CREATE CONSTRAINT IF NOT EXISTS ON (n:Vaccine) ASSERT exists(n.vaccine);
CREATE CONSTRAINT IF NOT EXISTS ON (n:Country) ASSERT exists(n.name);
