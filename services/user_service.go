//文件名称 /services/user_service.go业务逻辑代码

package services

import (
	"errors"
	datamodels2 "my.web/web/datamodels"
	repositories2 "my.web/web/repositories"
)

//UserService处理用户数据模型的CRUID操作，
//它取决于用户存储库的操作。
//这是将数据源与更高级别的组件分离。
//因此，不同的存储库类型可以使用相同的逻辑，而无需任何更改。
//它是一个接口，它在任何地方都被用作接口
//因为我们可能需要在将来更改或尝试实验性的不同域逻辑。
type UserService interface {
	GetAll() []datamodels2.User
	GetByID(id int64) (datamodels2.User, bool)
	GetByUsernameAndPassword(username, userPassword string) (datamodels2.User, bool)
	DeleteByID(id int64) bool
	Update(id int64, user datamodels2.User) (datamodels2.User, error)
	UpdatePassword(id int64, newPassword string) (datamodels2.User, error)
	UpdateUsername(id int64, newUsername string) (datamodels2.User, error)
	Create(userPassword string, user datamodels2.User) (datamodels2.User, error)
}

// NewUserService返回默认用户服务
func NewUserService(repo repositories2.UserRepository) UserService {
	return &UserServiceImplement{
		repo: repo,
	}
}

type UserServiceImplement struct {
	repo repositories2.UserRepository
}

// GetAll返回所有用户。
func (s *UserServiceImplement) GetAll() []datamodels2.User {
	return s.repo.SelectMany(func(_ datamodels2.User) bool {
		return true
	}, -1)
}

// GetByID根据其id返回用户。
func (s *UserServiceImplement) GetByID(id int64) (datamodels2.User, bool) {
	return s.repo.Select(func(m datamodels2.User) bool {
		return m.ID == id
	})
}

//获取yUsernameAndPassword根据用户名和密码返回用户，
//用于身份验证。
func (s *UserServiceImplement) GetByUsernameAndPassword(username, userPassword string) (datamodels2.User, bool) {
	if username == "" || userPassword == "" {
		return datamodels2.User{}, false
	}
	return s.repo.Select(func(m datamodels2.User) bool {
		if m.Username == username {
			hashed := m.HashedPassword
			if ok, _ := datamodels2.ValidatePassword(userPassword, hashed); ok {
				return true
			}
		}
		return false
	})
}

//更新现有用户的每个字段的更新，
//通过公共API使用是不安全的
//但是我们将在web  controllers/user_controller.go#PutBy上使用它
//为了向您展示它是如何工作的。
func (s *UserServiceImplement) Update(id int64, user datamodels2.User) (datamodels2.User, error) {
	user.ID = id
	return s.repo.InsertOrUpdate(user)
}

// UpdatePassword更新用户的密码。
func (s *UserServiceImplement) UpdatePassword(id int64, newPassword string) (datamodels2.User, error) {
	//更新用户并将其返回。
	hashed, err := datamodels2.GeneratePassword(newPassword)
	if err != nil {
		return datamodels2.User{}, err
	}
	return s.Update(id, datamodels2.User{
		HashedPassword: hashed,
	})
}

// UpdateUsername更新用户的用户名
func (s *UserServiceImplement) UpdateUsername(id int64, newUsername string) (datamodels2.User, error) {
	return s.Update(id, datamodels2.User{
		Username: newUsername,
	})
}

//创建插入新用户，
// userPassword是客户端类型的密码
//它将在插入我们的存储库之前进行哈希处理
func (s *UserServiceImplement) Create(userPassword string, user datamodels2.User) (datamodels2.User, error) {
	if user.ID > 0 || userPassword == "" || user.Firstname == "" || user.Username == "" {
		return datamodels2.User{}, errors.New("unable to create this user")
	}
	hashed, err := datamodels2.GeneratePassword(userPassword)
	if err != nil {
		return datamodels2.User{}, err
	}
	user.HashedPassword = hashed
	return s.repo.InsertOrUpdate(user)
}

// DeleteByID按其id删除用户。
//如果删除则返回true，否则返回false。
func (s *UserServiceImplement) DeleteByID(id int64) bool {
	return s.repo.Delete(func(m datamodels2.User) bool {
		return m.ID == id
	}, 1)
}
