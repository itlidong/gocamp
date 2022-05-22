// 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

// DAO module
if err == sql.ErrNoRows {
    return errors.Wrap(err, "some context inforamtion")
}

// Business Logic module
if errors.Is(err, sql.ErrNoRows} {
  // do something to handle the error
}
