import { type HTMLAttributes } from 'react';
import { cn } from '../lib/utils';
import { formatPhoneNumber, formatTimestamp } from '../lib/utils';
import Avatar from './Avatar';
import type { Contact } from '../types';

export interface ContactCardProps extends HTMLAttributes<HTMLDivElement> {
  contact: Contact;
  active?: boolean;
}

export default function ContactCard({ contact, active, className, ...props }: ContactCardProps) {
  const hasUnread = contact.unread_count > 0;

  return (
    <div
      className={cn(
        'flex items-center gap-3 bg-white rounded-lg p-4',
        'hover:bg-strawberry-50 cursor-pointer transition-colors',
        active && 'bg-strawberry-50',
        className
      )}
      {...props}
    >
      <Avatar name={contact.name} size="md" />

      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2">
          <h4 className={cn('font-medium text-neutral-800 truncate', hasUnread && 'font-semibold')}>
            {contact.name}
          </h4>
          {contact.last_message_at && (
            <span className="text-xs text-neutral-500 whitespace-nowrap">
              {formatTimestamp(contact.last_message_at)}
            </span>
          )}
        </div>
        <p className="text-sm text-neutral-500 truncate">{formatPhoneNumber(contact.phone_number)}</p>
      </div>

      {hasUnread && (
        <div className="px-2 py-0.5 bg-strawberry-500 text-white text-xs rounded-full min-w-[1.25rem] text-center">
          {contact.unread_count}
        </div>
      )}
    </div>
  );
}
