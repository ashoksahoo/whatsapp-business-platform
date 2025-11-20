import { type HTMLAttributes } from 'react';
import { cn } from '../lib/utils';
import { getInitials } from '../lib/utils';

export interface AvatarProps extends HTMLAttributes<HTMLDivElement> {
  name: string;
  src?: string;
  size?: 'sm' | 'md' | 'lg';
  online?: boolean;
}

export default function Avatar({ name, src, size = 'md', online, className, ...props }: AvatarProps) {
  const sizes = {
    sm: 'w-8 h-8 text-xs',
    md: 'w-12 h-12 text-base',
    lg: 'w-16 h-16 text-xl',
  };

  const onlineDotSizes = {
    sm: 'w-2 h-2',
    md: 'w-3 h-3',
    lg: 'w-4 h-4',
  };

  return (
    <div className={cn('relative', className)} {...props}>
      <div
        className={cn(
          'rounded-full flex items-center justify-center',
          'bg-strawberry-100 text-strawberry-700 font-semibold',
          sizes[size]
        )}
      >
        {src ? (
          <img src={src} alt={name} className="w-full h-full object-cover rounded-full" />
        ) : (
          <span>{getInitials(name)}</span>
        )}
      </div>
      {online && (
        <div
          className={cn(
            'absolute bottom-0 right-0 bg-leaf-500 border-2 border-white rounded-full',
            onlineDotSizes[size]
          )}
        />
      )}
    </div>
  );
}
