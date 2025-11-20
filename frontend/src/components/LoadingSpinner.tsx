import { type HTMLAttributes } from 'react';
import { cn } from '../lib/utils';

export interface LoadingSpinnerProps extends HTMLAttributes<HTMLDivElement> {
  size?: 'sm' | 'md' | 'lg';
}

export default function LoadingSpinner({ size = 'md', className, ...props }: LoadingSpinnerProps) {
  const sizes = {
    sm: 'w-4 h-4 border-2',
    md: 'w-8 h-8 border-4',
    lg: 'w-12 h-12 border-4',
  };

  return (
    <div
      className={cn(
        'border-neutral-200 border-t-strawberry-500 rounded-full animate-spin',
        sizes[size],
        className
      )}
      {...props}
    />
  );
}
