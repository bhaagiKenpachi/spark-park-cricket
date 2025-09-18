import { useDispatch, useSelector, TypedUseSelectorHook } from 'react-redux';
import type { RootState, AppDispatch } from './index';

export const useAppDispatch = () => {
  const dispatch = useDispatch<AppDispatch>();
  return dispatch;
};

export const useAppSelector: TypedUseSelectorHook<RootState> = selector => {
  const result = useSelector(selector);
  return result;
};
