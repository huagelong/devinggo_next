import { onUnmounted, ref } from 'vue';

import { useUserStore } from '@vben/stores';

import type { NotificationNewData, PresenceMemberData } from './pusher-types';
import { Events } from './pusher-types';
import { usePusher } from './use-pusher';

/**
 * Composable that subscribes to real-time notifications and online-user presence.
 * Automatically binds/unbinds on component mount/unmount.
 */
export function useRealtimeNotifications() {
  const { subscribePrivate, subscribePresence, bind, disconnect, state } =
    usePusher();

  const unreadCount = ref(0);
  const onlineUsers = ref<PresenceMemberData[]>([]);
  const latestNotification = ref<NotificationNewData | null>(null);

  const cleanups: (() => void)[] = [];

  /** Start listening for real-time events */
  function start() {
    const userStore = useUserStore();
    const userId =
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (userStore.userInfo as any)?.userId || (userStore.userInfo as any)?.id;

    if (!userId) {
      console.warn(
        '[Pusher] No user ID found, skipping real-time subscription',
      );
      return;
    }

    // 1. Private user channel — personal notifications
    const userChannel = subscribePrivate(`user-${userId}`);

    cleanups.push(
      bind(
        userChannel,
        Events.NOTIFICATION_NEW,
        (data: NotificationNewData) => {
          latestNotification.value = data;
          unreadCount.value++;
        },
      ),
    );

    cleanups.push(
      bind(userChannel, Events.NOTIFICATION_READ, () => {
        unreadCount.value = Math.max(0, unreadCount.value - 1);
      }),
    );

    // 2. Presence channel — online admin users
    const presenceChannel = subscribePresence('admin');

    cleanups.push(
      bind(
        presenceChannel,
        Events.SUBSCRIPTION_SUCCESS,
        (data: { members: PresenceMemberData[] }) => {
          onlineUsers.value = data.members || [];
        },
      ),
    );

    cleanups.push(
      bind(
        presenceChannel,
        Events.MEMBER_ADDED,
        (member: PresenceMemberData) => {
          onlineUsers.value = [...onlineUsers.value, member];
        },
      ),
    );

    cleanups.push(
      bind(
        presenceChannel,
        Events.MEMBER_REMOVED,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        (member: any) => {
          onlineUsers.value = onlineUsers.value.filter(
            (u) => String(u.id) !== String(member.id),
          );
        },
      ),
    );
  }

  /** Stop all subscriptions */
  function stop() {
    cleanups.forEach((fn) => fn());
    cleanups.length = 0;
    disconnect();
  }

  function setUnreadCount(count: number) {
    unreadCount.value = count;
  }

  function markRead(count = 1) {
    unreadCount.value = Math.max(0, unreadCount.value - count);
  }

  onUnmounted(() => {
    stop();
  });

  return {
    start,
    stop,
    unreadCount,
    onlineUsers,
    latestNotification,
    setUnreadCount,
    markRead,
    connectionState: state,
  };
}
