import { posthog } from './posthog';
import {
  PostHogEvent,
  SeriesEventProperties,
  SeriesPaginationProperties,
  MatchEventProperties,
  ScorecardEventProperties,
  BallEventProperties,
  OverEventProperties,
  WebSocketEventProperties,
  UserEventProperties,
  UserProperties,
} from '@/types/posthog';

// Helper to check if PostHog is available and enabled
const isPostHogAvailable = (): boolean => {
  // Disable all tracking in development
  if (process.env.NODE_ENV === 'development') {
    return false;
  }
  return typeof window !== 'undefined' && posthog !== null;
};

// Generic event tracking function
export const trackEvent = (
  eventName: PostHogEvent | string,
  properties?: Record<string, unknown>
): void => {
  if (!isPostHogAvailable()) return;

  try {
    posthog.capture(eventName, properties);
  } catch (error) {
    console.error('Failed to track event:', eventName, error);
  }
};

// Series Events
export const trackSeriesViewed = (properties: SeriesEventProperties): void => {
  trackEvent(PostHogEvent.SERIES_VIEWED, properties as unknown as Record<string, unknown>);
};

export const trackSeriesCreated = (properties: SeriesEventProperties): void => {
  trackEvent(PostHogEvent.SERIES_CREATED, properties as unknown as Record<string, unknown>);
};

export const trackSeriesEdited = (properties: SeriesEventProperties): void => {
  trackEvent(PostHogEvent.SERIES_EDITED, properties as unknown as Record<string, unknown>);
};

export const trackSeriesDeleted = (properties: SeriesEventProperties): void => {
  trackEvent(PostHogEvent.SERIES_DELETED, properties as unknown as Record<string, unknown>);
};

export const trackSeriesPaginationChanged = (
  properties: SeriesPaginationProperties
): void => {
  trackEvent(PostHogEvent.SERIES_PAGINATION_CHANGED, properties as unknown as Record<string, unknown>);
};

// Match Events
export const trackMatchViewed = (properties: MatchEventProperties): void => {
  trackEvent(PostHogEvent.MATCH_VIEWED, properties as unknown as Record<string, unknown>);
};

export const trackMatchCreated = (properties: MatchEventProperties): void => {
  trackEvent(PostHogEvent.MATCH_CREATED, properties as unknown as Record<string, unknown>);
};

export const trackMatchEdited = (properties: MatchEventProperties): void => {
  trackEvent(PostHogEvent.MATCH_EDITED, properties as unknown as Record<string, unknown>);
};

export const trackMatchDeleted = (properties: MatchEventProperties): void => {
  trackEvent(PostHogEvent.MATCH_DELETED, properties as unknown as Record<string, unknown>);
};

// Scorecard Events
export const trackScorecardViewed = (
  properties: ScorecardEventProperties
): void => {
  trackEvent(PostHogEvent.SCORECARD_VIEWED, properties as unknown as Record<string, unknown>);
};

export const trackLiveScoringStarted = (
  properties: ScorecardEventProperties
): void => {
  trackEvent(PostHogEvent.LIVE_SCORING_STARTED, properties as unknown as Record<string, unknown>);
};

export const trackBallAdded = (properties: BallEventProperties): void => {
  trackEvent(PostHogEvent.BALL_ADDED, properties as unknown as Record<string, unknown>);
};

export const trackInningsCompleted = (
  properties: ScorecardEventProperties
): void => {
  trackEvent(PostHogEvent.INNINGS_COMPLETED, properties as unknown as Record<string, unknown>);
};

export const trackMatchCompleted = (
  properties: ScorecardEventProperties
): void => {
  trackEvent(PostHogEvent.MATCH_COMPLETED, properties as unknown as Record<string, unknown>);
};

export const trackOverCompleted = (properties: OverEventProperties): void => {
  trackEvent(PostHogEvent.OVER_COMPLETED, properties as unknown as Record<string, unknown>);
};

// WebSocket Events
export const trackWebSocketConnected = (
  properties: WebSocketEventProperties
): void => {
  trackEvent(PostHogEvent.WEBSOCKET_CONNECTED, properties as unknown as Record<string, unknown>);
};

export const trackWebSocketDisconnected = (
  properties: WebSocketEventProperties
): void => {
  trackEvent(PostHogEvent.WEBSOCKET_DISCONNECTED, properties as unknown as Record<string, unknown>);
};

export const trackWebSocketError = (
  properties: WebSocketEventProperties
): void => {
  trackEvent(PostHogEvent.WEBSOCKET_ERROR, properties as unknown as Record<string, unknown>);
};

// User Events
export const trackUserLogin = (properties: UserEventProperties): void => {
  trackEvent(PostHogEvent.USER_LOGGED_IN, properties as unknown as Record<string, unknown>);
};

export const trackUserLogout = (properties: UserEventProperties): void => {
  trackEvent(PostHogEvent.USER_LOGGED_OUT, properties as unknown as Record<string, unknown>);
};

export const trackUserProfileUpdated = (
  properties: UserEventProperties
): void => {
  trackEvent(PostHogEvent.USER_PROFILE_UPDATED, properties as unknown as Record<string, unknown>);
};

// User Identification
export const identifyUser = (userId: string, properties?: UserProperties): void => {
  if (!isPostHogAvailable()) return;

  try {
    posthog.identify(userId, properties);
  } catch (error) {
    console.error('Failed to identify user:', error);
  }
};

export const resetUser = (): void => {
  if (!isPostHogAvailable()) return;

  try {
    posthog.reset();
  } catch (error) {
    console.error('Failed to reset user:', error);
  }
};

// Set user properties
export const setUserProperties = (properties: Partial<UserProperties>): void => {
  if (!isPostHogAvailable()) return;

  try {
    posthog.setPersonProperties(properties);
  } catch (error) {
    console.error('Failed to set user properties:', error);
  }
};

// Group analytics (for teams/organizations in future)
export const setGroup = (groupType: string, groupId: string): void => {
  if (!isPostHogAvailable()) return;

  try {
    posthog.group(groupType, groupId);
  } catch (error) {
    console.error('Failed to set group:', error);
  }
};

// Page view tracking (if manual tracking needed)
export const trackPageView = (path?: string): void => {
  if (!isPostHogAvailable()) return;

  try {
    posthog.capture('$pageview', { path: path || window.location.pathname });
  } catch (error) {
    console.error('Failed to track page view:', error);
  }
};

// Error Events
export const trackError = (properties: {
  error_type: string;
  error_message: string;
  error_stack?: string;
  component?: string;
  user_action?: string;
}): void => {
  trackEvent('error_occurred', properties);
};

export const trackApiError = (properties: {
  endpoint: string;
  status_code: number;
  error_message: string;
  request_duration?: number;
}): void => {
  trackEvent('api_error', properties);
};

// Performance Events
export const trackPerformance = (properties: {
  metric_name: string;
  metric_value: number;
  page?: string;
  component?: string;
}): void => {
  trackEvent('performance_metric', properties);
};

// User Engagement Events
export const trackFeatureUsage = (feature: string, metadata?: Record<string, unknown>): void => {
  trackEvent('feature_used', { feature, ...metadata });
};

export const trackUserAction = (action: string, metadata?: Record<string, unknown>): void => {
  trackEvent('user_action', { action, ...metadata });
};

export const trackTimeOnPage = (page: string, duration: number): void => {
  trackEvent('time_on_page', { page, duration });
};

// Funnel and Journey Tracking
export const trackFunnelStep = (
  funnelName: string,
  step: string,
  stepNumber: number,
  metadata?: Record<string, unknown>
): void => {
  trackEvent('funnel_step', {
    funnel_name: funnelName,
    step,
    step_number: stepNumber,
    ...metadata
  });
};

// Business Intelligence Events
export const trackConversion = (
  conversionType: string,
  value?: number,
  metadata?: Record<string, unknown>
): void => {
  trackEvent('conversion', {
    conversion_type: conversionType,
    value,
    ...metadata
  });
};

export const trackAchievement = (
  achievement: string,
  metadata?: Record<string, unknown>
): void => {
  trackEvent('achievement_unlocked', {
    achievement,
    ...metadata
  });
};

