package seeders

import (
	"time"

	"github.com/BangNopall/paskihub-be/domain/entity"
	"github.com/BangNopall/paskihub-be/domain/enums"
	"github.com/BangNopall/paskihub-be/pkg/bcrypt"
	"github.com/BangNopall/paskihub-be/pkg/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const demoPassword = "password123"

func DemoSeeder(db *gorm.DB) {
	hashedPassword, err := bcrypt.Bcrypt.Hash(demoPassword)
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[DEMO SEEDER] Failed to hash demo password")
	}

	now := time.Now()
	openDate := now.AddDate(0, -1, 0)
	closeDate := now.AddDate(0, 1, 0)
	compeDate := now.AddDate(0, 2, 0)
	pastDate := now.AddDate(0, -3, 0)

	users := buildDemoUsers(hashedPassword)
	institutions := buildDemoInstitutions()
	events := buildDemoEvents(openDate, closeDate, compeDate, pastDate)
	eventLevels := buildDemoEventLevels()
	wallets := buildDemoWallets()
	walletTransactions := buildDemoWalletTransactions()
	teams := buildDemoTeams()
	teamMembers := buildDemoTeamMembers()
	registrations := buildDemoRegistrations()
	judges := buildDemoJudges()
	scoreCategories := buildDemoScoreCategories()
	scoreSubCategories := buildDemoScoreSubCategories()
	gradeRules := buildDemoGradeRules()
	violationTypes := buildDemoViolationTypes()
	scores := buildDemoScores()
	teamViolations := buildDemoTeamViolations()
	eventAwards := buildDemoEventAwards()

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := insertDemoRecords(tx, users); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, institutions); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, events); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, eventLevels); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, wallets); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, walletTransactions); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, teams); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, teamMembers); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, registrations); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, judges); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, scoreCategories); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, scoreSubCategories); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, gradeRules); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, violationTypes); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, scores); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, teamViolations); err != nil {
			return err
		}
		if err := insertDemoRecords(tx, eventAwards); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatal(log.LogInfo{
			"error": err.Error(),
		}, "[DEMO SEEDER] Failed to seed demo data")
	}

	log.Info(log.LogInfo{
		"password": demoPassword,
	}, "[DEMO SEEDER] Demo data seeded successfully")
}

func insertDemoRecords[T any](db *gorm.DB, records []T) error {
	if len(records) == 0 {
		return nil
	}

	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&records).Error
}

func demoID(name string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte("paskihub-demo:"+name))
}

func buildDemoUsers(hashedPassword string) []entity.User {
	return []entity.User{
		{
			Id:              demoID("user-admin"),
			Email:           "admin.demo@paskihub.local",
			Password:        hashedPassword,
			Role:            enums.Admin,
			EmailIsVerified: true,
		},
		{
			Id:              demoID("user-organizer"),
			Email:           "organizer.demo@paskihub.local",
			Password:        hashedPassword,
			Role:            enums.Organizer,
			EmailIsVerified: true,
		},
		{
			Id:              demoID("user-peserta"),
			Email:           "peserta.demo@paskihub.local",
			Password:        hashedPassword,
			Role:            enums.Peserta,
			EmailIsVerified: true,
		},
	}
}

func buildDemoInstitutions() []entity.Institution {
	return []entity.Institution{
		{
			Id:              demoID("institution-sma-demo"),
			UserId:          demoID("user-peserta"),
			Name:            "SMA Paskihub Demo",
			Address:         "Jl. Merdeka No. 17, Bandung",
			InstitutionType: enums.SMA,
			NamePj:          "Raka Demo",
			NoWaPj:          "+6281234567890",
		},
	}
}

func buildDemoEvents(openDate, closeDate, compeDate, pastDate time.Time) []entity.Event {
	return []entity.Event{
		{
			Id:             demoID("event-open"),
			UserId:         demoID("user-organizer"),
			Name:           "Paskibra Open Demo 2026",
			Description:    "Event demo untuk menguji alur pendaftaran, wallet, dan penilaian.",
			OpenDate:       openDate,
			CloseDate:      closeDate,
			CompeDate:      compeDate,
			Organizer:      "Paskihub Demo Organizer",
			Location:       "GOR Demo Bandung",
			MinTeamMembers: 8,
			MaxTeamMembers: 18,
			Status:         string(enums.Open),
			NamePj:         "Nadia Demo",
			NoWaPj:         "+6281111111111",
			BankName:       "BCA",
			BankNumber:     "1234567890",
			WaGroup:        "https://chat.whatsapp.com/demo-open",
		},
		{
			Id:             demoID("event-draft"),
			UserId:         demoID("user-organizer"),
			Name:           "Lomba Purna Demo 2026",
			Description:    "Event draft untuk menguji data yang belum dibuka.",
			OpenDate:       nowAt(openDate, 1),
			CloseDate:      nowAt(closeDate, 1),
			CompeDate:      nowAt(compeDate, 1),
			Organizer:      "Paskihub Demo Organizer",
			Location:       "Lapangan Demo Jakarta",
			MinTeamMembers: 6,
			MaxTeamMembers: 14,
			Status:         string(enums.Draft),
			NamePj:         "Nadia Demo",
			NoWaPj:         "+6281111111111",
			BankName:       "Mandiri",
			BankNumber:     "9876543210",
			WaGroup:        "https://chat.whatsapp.com/demo-draft",
		},
		{
			Id:             demoID("event-closed"),
			UserId:         demoID("user-organizer"),
			Name:           "Kejuaraan Arsip Demo 2025",
			Description:    "Event selesai untuk menguji riwayat dan rekap.",
			OpenDate:       pastDate.AddDate(0, -1, 0),
			CloseDate:      pastDate.AddDate(0, 0, 14),
			CompeDate:      pastDate.AddDate(0, 0, 21),
			Organizer:      "Paskihub Demo Organizer",
			Location:       "Alun-alun Demo",
			MinTeamMembers: 8,
			MaxTeamMembers: 18,
			Status:         string(enums.Closed),
			NamePj:         "Nadia Demo",
			NoWaPj:         "+6281111111111",
			BankName:       "BNI",
			BankNumber:     "1122334455",
			WaGroup:        "https://chat.whatsapp.com/demo-closed",
		},
	}
}

func buildDemoEventLevels() []entity.EventLevel {
	return []entity.EventLevel{
		{
			Id:       demoID("level-open-smp"),
			EventId:  demoID("event-open"),
			Name:     "SMP",
			RegisFee: "250000",
			DpFee:    "100000",
		},
		{
			Id:       demoID("level-open-sma"),
			EventId:  demoID("event-open"),
			Name:     "SMA/SMK",
			RegisFee: "350000",
			DpFee:    "150000",
		},
		{
			Id:       demoID("level-draft-umum"),
			EventId:  demoID("event-draft"),
			Name:     "UMUM",
			RegisFee: "300000",
			DpFee:    "125000",
		},
		{
			Id:               demoID("level-closed-sma"),
			EventId:          demoID("event-closed"),
			Name:             "SMA/SMK",
			RegisFee:         "300000",
			DpFee:            "100000",
			IsScorePublished: true,
		},
	}
}

func buildDemoWallets() []entity.Wallet {
	return []entity.Wallet{
		{Id: demoID("wallet-open"), EventId: demoID("event-open"), Saldo: 750000},
		{Id: demoID("wallet-draft"), EventId: demoID("event-draft"), Saldo: 100000},
		{Id: demoID("wallet-closed"), EventId: demoID("event-closed"), Saldo: 300000},
	}
}

func buildDemoWalletTransactions() []entity.WalletTransaction {
	return []entity.WalletTransaction{
		{
			Id:       demoID("wallet-transaction-approved"),
			WalletId: demoID("wallet-open"),
			Type:     enums.TopUp,
			Amount:   500000,
			Status:   enums.Approve,
		},
		{
			Id:       demoID("wallet-transaction-pending"),
			WalletId: demoID("wallet-open"),
			Type:     enums.TopUp,
			Amount:   250000,
			Status:   enums.Pending,
		},
		{
			Id:              demoID("wallet-transaction-rejected"),
			WalletId:        demoID("wallet-draft"),
			Type:            enums.TopUp,
			Amount:          100000,
			Status:          enums.TSRejected,
			RejectionReason: "Bukti transfer tidak terbaca",
		},
	}
}

func buildDemoTeams() []entity.Team {
	return []entity.Team{
		{Id: demoID("team-garuda"), InstiId: demoID("institution-sma-demo"), Name: "Garuda Demo", Pelatih: "Pak Dimas Demo"},
		{Id: demoID("team-melati"), InstiId: demoID("institution-sma-demo"), Name: "Melati Demo", Pelatih: "Bu Sari Demo"},
		{Id: demoID("team-cakra"), InstiId: demoID("institution-sma-demo"), Name: "Cakra Demo", Pelatih: "Pak Bima Demo"},
		{Id: demoID("team-nusa"), InstiId: demoID("institution-sma-demo"), Name: "Nusa Demo", Pelatih: "Bu Rani Demo"},
	}
}

func buildDemoTeamMembers() []entity.TeamMember {
	teams := []string{"garuda", "melati", "cakra", "nusa"}
	members := make([]entity.TeamMember, 0, len(teams)*4)

	for _, team := range teams {
		teamID := demoID("team-" + team)
		members = append(members,
			entity.TeamMember{Id: demoID("member-" + team + "-danpas"), TeamId: teamID, FullName: "Danpas " + team, Role: enums.Danpas},
			entity.TeamMember{Id: demoID("member-" + team + "-pasukan-1"), TeamId: teamID, FullName: "Pasukan 1 " + team, Role: enums.Pasukan},
			entity.TeamMember{Id: demoID("member-" + team + "-official"), TeamId: teamID, FullName: "Official " + team, Role: enums.Official},
			entity.TeamMember{Id: demoID("member-" + team + "-pelatih"), TeamId: teamID, FullName: "Pelatih " + team, Role: enums.Pelatih},
		)
	}

	return members
}

func buildDemoRegistrations() []entity.Registration {
	return []entity.Registration{
		{
			Id:               demoID("registration-full-paid"),
			TeamId:           demoID("team-garuda"),
			EventLevelId:     demoID("level-open-sma"),
			PaymentType:      "LUNAS",
			PaymentStatus:    enums.FullPaid,
			AssessmentStatus: enums.AssessCompleted,
			PaymentProofPath: "public/demo/payment-full-paid.jpg",
		},
		{
			Id:               demoID("registration-waiting"),
			TeamId:           demoID("team-melati"),
			EventLevelId:     demoID("level-open-smp"),
			PaymentType:      "DP",
			PaymentStatus:    enums.Waiting,
			AssessmentStatus: enums.AssessPending,
			PaymentProofPath: "public/demo/payment-waiting.jpg",
		},
		{
			Id:               demoID("registration-dp-paid"),
			TeamId:           demoID("team-cakra"),
			EventLevelId:     demoID("level-open-sma"),
			PaymentType:      "DP",
			PaymentStatus:    enums.DPPaid,
			AssessmentStatus: enums.AssessInProgress,
			PaymentProofPath: "public/demo/payment-dp-paid.jpg",
		},
		{
			Id:               demoID("registration-rejected"),
			TeamId:           demoID("team-nusa"),
			EventLevelId:     demoID("level-closed-sma"),
			PaymentType:      "LUNAS",
			PaymentStatus:    enums.Rejected,
			AssessmentStatus: enums.AssessPending,
			PaymentProofPath: "public/demo/payment-rejected.jpg",
			RejectionReason:  "Bukti pembayaran tidak sesuai nominal",
		},
	}
}

func buildDemoJudges() []entity.Judge {
	return []entity.Judge{
		{ID: demoID("judge-1"), EventID: demoID("event-open"), Name: "Juri Demo 1"},
		{ID: demoID("judge-2"), EventID: demoID("event-open"), Name: "Juri Demo 2"},
	}
}

func buildDemoScoreCategories() []entity.ScoreCategory {
	return []entity.ScoreCategory{
		{ID: demoID("score-category-formasi"), EventID: demoID("event-open"), EventLevelID: demoID("level-open-sma"), Name: "Formasi"},
		{ID: demoID("score-category-gerakan"), EventID: demoID("event-open"), EventLevelID: demoID("level-open-sma"), Name: "Gerakan"},
	}
}

func buildDemoScoreSubCategories() []entity.ScoreSubCategory {
	return []entity.ScoreSubCategory{
		{ID: demoID("score-subcategory-kekompakan"), ScoreCategoriesID: demoID("score-category-formasi"), Name: "Kekompakan", MaxScore: 100},
		{ID: demoID("score-subcategory-kerapian"), ScoreCategoriesID: demoID("score-category-formasi"), Name: "Kerapian", MaxScore: 100},
		{ID: demoID("score-subcategory-ketegasan"), ScoreCategoriesID: demoID("score-category-gerakan"), Name: "Ketegasan", MaxScore: 100},
	}
}

func buildDemoGradeRules() []entity.GradeRule {
	return []entity.GradeRule{
		{ID: demoID("grade-kurang"), EventID: demoID("event-open"), EventLevelID: demoID("level-open-sma"), GradeName: "Kurang", MinScore: 0, MaxScore: 59},
		{ID: demoID("grade-cukup"), EventID: demoID("event-open"), EventLevelID: demoID("level-open-sma"), GradeName: "Cukup", MinScore: 60, MaxScore: 74},
		{ID: demoID("grade-baik"), EventID: demoID("event-open"), EventLevelID: demoID("level-open-sma"), GradeName: "Baik", MinScore: 75, MaxScore: 89},
		{ID: demoID("grade-sangat-baik"), EventID: demoID("event-open"), EventLevelID: demoID("level-open-sma"), GradeName: "Sangat Baik", MinScore: 90, MaxScore: 100},
	}
}

func buildDemoViolationTypes() []entity.ViolationType {
	return []entity.ViolationType{
		{ID: demoID("violation-seragam"), EventID: demoID("event-open"), EventLevelID: demoID("level-open-sma"), Name: "Atribut tidak lengkap", Point: 5},
		{ID: demoID("violation-waktu"), EventID: demoID("event-open"), EventLevelID: demoID("level-open-sma"), Name: "Melebihi batas waktu", Point: 10},
	}
}

func buildDemoScores() []entity.Score {
	return []entity.Score{
		{ID: demoID("score-1"), RegisID: demoID("registration-full-paid"), JudgesID: demoID("judge-1"), SubCategoryID: demoID("score-subcategory-kekompakan"), ScoreValue: 88, Grade: "Baik"},
		{ID: demoID("score-2"), RegisID: demoID("registration-full-paid"), JudgesID: demoID("judge-1"), SubCategoryID: demoID("score-subcategory-kerapian"), ScoreValue: 92, Grade: "Sangat Baik"},
		{ID: demoID("score-3"), RegisID: demoID("registration-full-paid"), JudgesID: demoID("judge-2"), SubCategoryID: demoID("score-subcategory-kekompakan"), ScoreValue: 84, Grade: "Baik"},
		{ID: demoID("score-4"), RegisID: demoID("registration-full-paid"), JudgesID: demoID("judge-2"), SubCategoryID: demoID("score-subcategory-ketegasan"), ScoreValue: 90, Grade: "Sangat Baik"},
	}
}

func buildDemoTeamViolations() []entity.TeamViolation {
	return []entity.TeamViolation{
		{ID: demoID("team-violation-1"), RegisID: demoID("registration-full-paid"), JudgesID: demoID("judge-1"), ViolationTypeID: demoID("violation-seragam")},
	}
}

func buildDemoEventAwards() []entity.EventAward {
	return []entity.EventAward{
		{ID: demoID("event-award-juara-utama"), EventID: demoID("event-open"), Name: "Juara Utama", LimitRank: 3},
		{ID: demoID("event-award-danton-terbaik"), EventID: demoID("event-open"), Name: "Danton Terbaik", LimitRank: 1},
	}
}

func nowAt(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}
