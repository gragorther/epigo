package gormdb_test

func (s *DBTestSuite) TestCheckIfGroupExistsByID() {
	s.Run("user doesn't exist", func() {
		exists, err := s.repo.CheckIfGroupExistsByID(s.ctx, 2)
		s.Require().NoError(err)
		s.False(exists, "user shouldn't exist")
	})
}
