import { type HTMLAttributes } from 'react';
import { cn } from '../lib/utils';

export interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: 'success' | 'error' | 'warning' | 'info' | 'neutral';
  icon?: React.ReactNode;
}

export default function Badge({ variant = 'neutral', icon, className, children, ...props }: BadgeProps) {
  const variants = {
    success: 'bg-leaf-100 text-leaf-700',
    error: 'bg-strawberry-100 text-strawberry-700',
    warning: 'bg-yellow-100 text-yellow-700',
    info: 'bg-blue-100 text-blue-700',
    neutral: 'bg-neutral-100 text-neutral-700',
  };

  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 px-2 py-1 text-xs font-medium rounded-full',
        variants[variant],
        className
      )}
      {...props}
    >
      {icon}
      {children}
    </span>
  );
}
