package generator

import (
	"context"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/syth0le/social-network/internal/model"
	"github.com/syth0le/social-network/internal/storage"
	"github.com/syth0le/social-network/internal/utils"
)

const (
	batchLengthDivider       = 1000
	batchLoggerWriterDivider = batchLengthDivider * 10
)

type Service interface {
	BatchGenerateUsers(ctx context.Context) error
}

type ServiceImpl struct {
	Logger   *zap.Logger
	Storage  storage.Storage
	DataFile string
}

func (s *ServiceImpl) BatchGenerateUsers(ctx context.Context) error {
	// TODO: make atomic transaction
	usersList := s.readUsersData()
	var usersListForBatchCreate []*model.UserRegister
	counter := 0
	for _, user := range usersList {
		usersListForBatchCreate = append(usersListForBatchCreate, user)

		err := user.Validate()
		if err != nil {
			return fmt.Errorf("user validate: %w", err)
		}

		if len(usersListForBatchCreate) == batchLengthDivider {
			err := s.Storage.User().BatchCreateUser(ctx, usersListForBatchCreate)
			if err != nil {
				return fmt.Errorf("batch create user: %w", err)
			}

			counter += batchLengthDivider
			if counter%batchLoggerWriterDivider == 0 {
				s.Logger.Sugar().Infof("created: %d", counter)
			}
			usersListForBatchCreate = nil
		}
	}

	if len(usersListForBatchCreate) != 0 {
		err := s.Storage.User().BatchCreateUser(ctx, usersListForBatchCreate)
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		counter += len(usersListForBatchCreate)
		s.Logger.Sugar().Infof("created: %d", counter)
	}

	return nil
}

func (s *ServiceImpl) readUsersData() []*model.UserRegister {
	s.Logger.Info("starting read data...")

	file, err := os.Open(s.DataFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	counter := 0
	var usersList []*model.UserRegister
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			fmt.Println(err)
			break
		}
		counter += 1

		fullName := strings.Split(record[0], " ")
		id := utils.GenerateUUID()
		age, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Println(err)
			break
		}

		usersList = append(usersList, &model.UserRegister{
			ID:             id,
			Username:       id,
			HashedPassword: "password",
			FirstName:      fullName[0],
			SecondName:     fullName[1],
			Sex:            "",
			Birthdate:      time.Time{}.AddDate(time.Now().Year()-age, rand.Intn(12), rand.Intn(28)),
			Biography:      "",
			City:           record[2],
		})
	}
	s.Logger.Sugar().Infof("reading data finished. %d rows read", counter)
	return usersList
}
