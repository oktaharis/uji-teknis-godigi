package database

import "gorm.io/gorm"

// karna ga bisa pakai gorm saya pakai DDL
func EnsureCompanyTables(db *gorm.DB) error {
	// LEADS
	if err := db.Exec(`
CREATE TABLE IF NOT EXISTS leads (
  lead_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  company_name VARCHAR(255) NOT NULL,
  contact_name VARCHAR(100) NOT NULL,
  email VARCHAR(255) NOT NULL,
  phone VARCHAR(30),
  source VARCHAR(50),
  industry VARCHAR(50),
  region VARCHAR(50),
  sales_rep VARCHAR(50),
  status VARCHAR(30),
  notes TEXT,
  PRIMARY KEY (lead_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`).Error; err != nil {
		return err
	}

	// DEALS (dengan FK ke leads)
	if err := db.Exec(`
CREATE TABLE IF NOT EXISTS deals (
  deal_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  lead_id BIGINT UNSIGNED NOT NULL,
  deal_name VARCHAR(120),
  amount_idr BIGINT NOT NULL,
  currency CHAR(3) NOT NULL DEFAULT 'IDR',
  term_months INT NOT NULL DEFAULT 12,
  stage VARCHAR(20) NOT NULL,
  closed_at TIMESTAMP NOT NULL,
  PRIMARY KEY (deal_id),
  CONSTRAINT fk_deals_lead FOREIGN KEY (lead_id)
    REFERENCES leads(lead_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`).Error; err != nil {
		return err
	}

	return nil
}
