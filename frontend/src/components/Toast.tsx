import { useEffect } from 'react';
import { X, Check, AlertCircle } from 'lucide-react';
import { cn } from '../lib/utils';

export interface ToastProps {
  id: string;
  type?: 'success' | 'error' | 'info';
  title: string;
  message?: string;
  duration?: number;
  onClose: (id: string) => void;
}

export default function Toast({ id, type = 'info', title, message, duration = 5000, onClose }: ToastProps) {
  useEffect(() => {
    if (duration > 0) {
      const timer = setTimeout(() => onClose(id), duration);
      return () => clearTimeout(timer);
    }
  }, [id, duration, onClose]);

  const icons = {
    success: <Check className="w-3 h-3 text-leaf-600" />,
    error: <AlertCircle className="w-3 h-3 text-strawberry-600" />,
    info: <AlertCircle className="w-3 h-3 text-blue-600" />,
  };

  const borderColors = {
    success: 'border-leaf-500',
    error: 'border-strawberry-500',
    info: 'border-blue-500',
  };

  const bgColors = {
    success: 'bg-leaf-100',
    error: 'bg-strawberry-100',
    info: 'bg-blue-100',
  };

  return (
    <div
      className={cn(
        'bg-white rounded-lg shadow-lg p-4',
        'flex items-start gap-3 max-w-sm',
        'border-l-4',
        borderColors[type],
        'animate-slideInRight'
      )}
    >
      <div className={cn('w-5 h-5 rounded-full flex items-center justify-center', bgColors[type])}>
        {icons[type]}
      </div>
      <div className="flex-1">
        <h4 className="text-sm font-semibold text-neutral-800">{title}</h4>
        {message && <p className="text-sm text-neutral-600 mt-0.5">{message}</p>}
      </div>
      <button
        onClick={() => onClose(id)}
        className="text-neutral-400 hover:text-neutral-600 transition-colors"
      >
        <X className="w-4 h-4" />
      </button>
    </div>
  );
}
