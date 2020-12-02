#作业

	我们在数据库操作的时候，比如dao层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层，为什么，应该怎么作请写出代码？

##解答

	当 dao 层遇到 sql.ErrNoRows 时，代表未查询到对应数据，这个 error 不应该抛给上层，dao 层需要处理这个错误，并向上层抛出自定义错误。因为这个错误是在对数据库进行处理时的错误，上层的错误处理不应该和数据库的错误产生关系。

##伪代码

	package dao
	import (
	    "database/sql"
	    "github.com/pkg/errors"
	    merrors "errors"
	)

	var NotFound = merrors.New("not found")

	type Dao struct {
	    db *sql.DB
	}

	func New() *Dao {
	    return &Dao{
		    db: &sql.DB{},
	    }
	}

	func (d *Dao) Find() (*User, error) {
		user, err := d.db.Query("select RowOne from User;")
		if err != nil {
			if merrors.Is(err, sql.ErrNoRows) {
				return nil, NotFound
			}
			return nil, errors.Wrap(err, "find error")
		}
		return &user, err
	}
