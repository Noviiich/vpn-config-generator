package service

// func (s *VPNService) StatusSubscribtion(ctx context.Context, username string, chatID int) (st string, err error) {
// 	defer func() { err = e.WrapIfErr("can't get status subscription", err) }()

// 	user, err := s.GetUser(ctx, chatID, username)
// 	if err != nil {
// 		return "", err
// 	}

// 	if user.SubscriptionActive {
// 		remaining := time.Until(user.SubscriptionExpiry)
// 		days := int(remaining.Hours()) / 24
// 		hours := int(remaining.Hours()) % 24

// 		msg := fmt.Sprintf(`Ваша подписка активна!!!
// Она истекает через %d дней, %d часов`, days, hours)
// 		return msg, nil
// 	}

// 	return `У вас не подписки. Не расстраивайтесь, вы все еще можете ее оформить.
// Для этого выполните /subscribe`, nil
// }

// func (s *VPNService) UpdateSubscription(ctx context.Context, chatID int) (err error) {
// 	defer func() { err = e.WrapIfErr("can't update subscription", err) }()

// 	user, err := s.repo.GetUser(ctx, chatID)
// 	if err != nil {
// 		return err
// 	}

// 	user.SubscriptionActive = true
// 	user.SubscriptionExpiry = time.Now().Add(30 * 24 * time.Hour)

// 	return s.repo.UpdateUser(ctx, user)
// }
