package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationService(t *testing.T) {
	testApp := SetupTestApp()
	defer TeardownTestApp(testApp)

	notificationService := testApp.NotificationService

	user1, _ := CreateTestUser(t, testApp, "noti_user1@test.com", "user")
	user2, _ := CreateTestUser(t, testApp, "noti_user2@test.com", "user")

	t.Run("Create and Get Notifications", func(t *testing.T) {
		_, err := notificationService.CreateNotification(user1.ID, "Test notification 1 for user 1", "/link1")
		require.NoError(t, err)
		_, err = notificationService.CreateNotification(user1.ID, "Test notification 2 for user 1", "/link2")
		require.NoError(t, err)
		_, err = notificationService.CreateNotification(user2.ID, "Test notification for user 2", "/link3")
		require.NoError(t, err)

		// Get all notifications for user1 (the service layer doesn't filter by read status)
		notifications, err := notificationService.GetUserNotifications(user1.ID)
		require.NoError(t, err)
		assert.Len(t, notifications, 2, "User1 should have 2 notifications")
		assert.Equal(t, "Test notification 1 for user 1", notifications[0].Message)

		notifications, err = notificationService.GetUserNotifications(user2.ID)
		require.NoError(t, err)
		assert.Len(t, notifications, 1, "User2 should have 1 notification")
	})

	t.Run("Mark as Read", func(t *testing.T) {
		user3, _ := CreateTestUser(t, testApp, "noti_user3@test.com", "user")
		noti, err := notificationService.CreateNotification(user3.ID, "Notification to be read", "/link4")
		require.NoError(t, err)

		// Mark as read
		err = notificationService.MarkNotificationAsRead(noti.ID, user3.ID)
		require.NoError(t, err)

		// Verify it's now read
		updatedNoti, err := notificationService.GetNotificationByID(noti.ID)
		require.NoError(t, err)
		assert.True(t, updatedNoti.IsRead, "Notification should be marked as read")
	})

	t.Run("Mark All as Read", func(t *testing.T) {
		user4, _ := CreateTestUser(t, testApp, "noti_user4@test.com", "user")
		for i := 0; i < 3; i++ {
			_, err := notificationService.CreateNotification(user4.ID, fmt.Sprintf("Noti %d", i), "/link")
			require.NoError(t, err)
		}

		// Mark all as read
		err := notificationService.MarkAllUserNotificationsAsRead(user4.ID)
		require.NoError(t, err)

		// Verify they are all read
		notifications, err := notificationService.GetUserNotifications(user4.ID)
		require.NoError(t, err)
		for _, n := range notifications {
			assert.True(t, n.IsRead, "All notifications should be marked as read")
		}
	})
}
