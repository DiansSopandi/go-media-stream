package middlewares

import (
	"github.com/DiansSopandi/media_stream/db"
	"github.com/gofiber/fiber/v2"
)

const UserServiceCtxKey = "userServiceWithTx"
const TxContextKey = "tx"

func WithTransaction(handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dbConn := db.InitDatabase()
		tx, err := dbConn.Begin() // start transaction commit rollback
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to start transaction")
		}

		// repo, err := repository.NewUserRepository(true) // true = start TX
		// if err != nil {
		// 	return fiber.NewError(fiber.StatusInternalServerError, "Failed to start transaction")
		// }

		// Inject service dengan repo yg punya TX ke context
		// userService := service.NewUserService(repo)
		// c.Locals(UserServiceCtxKey, userService)
		// Simpan TX ke dalam context
		c.Locals(TxContextKey, tx)

		// Jalankan handler utama
		err = handler(c)

		if err != nil {
			// Rollback jika ada error
			db.RollbackOnError(tx, err)
			// return fmt.Errorf("transaction rolled back due to: %v", err)
			return err
			// return fiber.NewError(
			// 	fiber.StatusInternalServerError,
			// 	fmt.Sprintf("transaction rolled back due to: %v", err),
			// )
			// return pkg.ResponseApiErrorInternalServer(c, fmt.Sprintf("transaction rolled back due to: %v", err))
		}

		// Commit transaksi jika sukses
		if tx != nil {
			if err := tx.Commit(); err != nil {
				// return fmt.Errorf("failed to commit transaction: %v", err)
				// only return error to cover error handler
				return err
				// return fiber.NewError(
				// 	fiber.StatusInternalServerError,
				// 	fmt.Sprintf("Failed to commit transaction: %v", err),
				// )
				// return pkg.ResponseApiErrorInternalServer(c, fmt.Sprintf("Failed to commit transaction: %v", err))
			}
		}

		return nil
	}
}
