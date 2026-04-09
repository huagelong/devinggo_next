import Pusher from 'pusher-js';
import type { Channel } from 'pusher-js';
import { readonly, ref } from 'vue';

import { useAccessStore } from '@vben/stores';

export interface PusherConfig {
  appKey: string;
  wsHost?: string;
  wsPort?: number;
  wssPort?: number;
  forceTLS?: boolean;
  authEndpoint?: string;
  cluster?: string;
}

// Default config — matches hack/config.yaml pusher section
const defaultConfig: PusherConfig = {
  appKey: 'devinggo-app-key',
  wsHost: window.location.hostname,
  wsPort: 8070,
  wssPort: 8070,
  forceTLS: window.location.protocol === 'https:',
  authEndpoint: '/system/pusher/auth',
  cluster: 'local',
};

let pusherInstance: Pusher | null = null;
const connectionState = ref<string>('disconnected');

/**
 * Initialize or return the singleton Pusher client.
 */
export function usePusher(config?: Partial<PusherConfig>) {
  const merged = { ...defaultConfig, ...config };

  function getInstance(): Pusher {
    if (pusherInstance) return pusherInstance;

    const accessStore = useAccessStore();
    const token = accessStore.accessToken ?? '';

    pusherInstance = new Pusher(merged.appKey, {
      cluster: merged.cluster || 'local',
      wsHost: merged.wsHost,
      wsPort: merged.wsPort,
      wssPort: merged.wssPort,
      forceTLS: merged.forceTLS,
      enabledTransports: ['ws', 'wss'],
      authEndpoint: merged.authEndpoint,
      auth: {
        headers: {
          Authorization: token,
        },
      },
    });

    // Connection state monitoring
    pusherInstance.connection.bind(
      'state_change',
      (states: { current: string; previous: string }) => {
        connectionState.value = states.current;
        console.debug(
          '[Pusher] state:',
          states.previous,
          '→',
          states.current,
        );
      },
    );

    pusherInstance.connection.bind('error', (err: Error) => {
      console.error('[Pusher] connection error:', err);
    });

    return pusherInstance;
  }

  /** Subscribe to a public channel */
  function subscribe(channelName: string): Channel {
    return getInstance().subscribe(channelName);
  }

  /** Subscribe to a private channel (requires auth) */
  function subscribePrivate(channelName: string): Channel {
    const name = channelName.startsWith('private-')
      ? channelName
      : `private-${channelName}`;
    return getInstance().subscribe(name);
  }

  /** Subscribe to a presence channel (requires auth, tracks users) */
  function subscribePresence(channelName: string): Channel {
    const name = channelName.startsWith('presence-')
      ? channelName
      : `presence-${channelName}`;
    return getInstance().subscribe(name);
  }

  /** Unsubscribe from a channel */
  function unsubscribe(channelName: string): void {
    pusherInstance?.unsubscribe(channelName);
  }

  /** Bind a callback to a channel event, returns unbind fn */
  function bind(
    channel: Channel,
    event: string,
    callback: (data: any) => void,
  ) {
    channel.bind(event, callback);
    return () => channel.unbind(event, callback);
  }

  /** Disconnect and clean up */
  function disconnect(): void {
    if (pusherInstance) {
      pusherInstance.disconnect();
      pusherInstance = null;
      connectionState.value = 'disconnected';
    }
  }

  return {
    getInstance,
    subscribe,
    subscribePrivate,
    subscribePresence,
    unsubscribe,
    bind,
    disconnect,
    state: readonly(connectionState),
  };
}
