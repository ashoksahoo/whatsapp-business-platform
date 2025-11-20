import { type HTMLAttributes } from 'react';
import { Check, CheckCheck } from 'lucide-react';
import { cn } from '../lib/utils';
import { formatTimestamp } from '../lib/utils';
import type { Message } from '../types';

export interface MessageBubbleProps extends HTMLAttributes<HTMLDivElement> {
  message: Message;
}

export default function MessageBubble({ message, className, ...props }: MessageBubbleProps) {
  const isOutbound = message.direction === 'outbound';

  const getStatusIcon = () => {
    if (message.direction === 'inbound') return null;

    switch (message.status) {
      case 'sent':
        return <Check className="w-4 h-4 text-neutral-500" />;
      case 'delivered':
        return <CheckCheck className="w-4 h-4 text-neutral-500" />;
      case 'read':
        return <CheckCheck className="w-4 h-4 text-leaf-500" />;
      case 'failed':
        return <span className="text-xs text-strawberry-600">Failed</span>;
      default:
        return null;
    }
  };

  return (
    <div className={cn('flex mb-4', isOutbound ? 'justify-end' : 'justify-start', className)} {...props}>
      <div className="max-w-[70%]">
        <div
          className={cn(
            'rounded-2xl px-4 py-2 shadow-sm',
            isOutbound ? 'bg-leaf-100 rounded-tr-sm' : 'bg-white rounded-tl-sm'
          )}
        >
          {message.media_url && (
            <div className="mb-2">
              {message.message_type === 'image' ? (
                <img
                  src={message.media_url}
                  alt="Media"
                  className="rounded-lg max-w-full h-auto"
                />
              ) : (
                <a
                  href={message.media_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-strawberry-600 hover:underline text-sm"
                >
                  View {message.message_type}
                </a>
              )}
            </div>
          )}
          <p className="text-neutral-800 whitespace-pre-wrap break-words">{message.content}</p>
        </div>
        <div className={cn('flex items-center gap-2 mt-1 px-2', isOutbound ? 'justify-end' : 'justify-start')}>
          <span className="text-xs text-neutral-500">{formatTimestamp(message.timestamp)}</span>
          {getStatusIcon()}
        </div>
      </div>
    </div>
  );
}
