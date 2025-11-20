import { type ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { MessageSquare, Users, FileText, Settings } from 'lucide-react';
import { cn } from '../lib/utils';

export interface MainLayoutProps {
  children: ReactNode;
}

const navigation = [
  { name: 'Messages', href: '/messages', icon: MessageSquare },
  { name: 'Contacts', href: '/contacts', icon: Users },
  { name: 'Templates', href: '/templates', icon: FileText },
  { name: 'Settings', href: '/settings', icon: Settings },
];

export default function MainLayout({ children }: MainLayoutProps) {
  const location = useLocation();

  return (
    <div className="h-screen flex flex-col bg-neutral-50">
      {/* Top Nav */}
      <nav className="h-16 bg-white border-b border-neutral-200 px-6 flex items-center justify-between flex-shrink-0">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-strawberry-500 rounded-lg flex items-center justify-center">
            <span className="text-white font-bold text-lg">🍓</span>
          </div>
          <h1 className="text-xl font-bold text-neutral-900">Vibecoded WA</h1>
        </div>
        <div className="flex items-center gap-4">
          <span className="text-sm text-neutral-600">Admin</span>
        </div>
      </nav>

      <div className="flex-1 flex overflow-hidden">
        {/* Sidebar */}
        <aside className="w-64 bg-white border-r border-neutral-200 flex-shrink-0">
          <nav className="p-4 space-y-1">
            {navigation.map((item) => {
              const isActive = location.pathname.startsWith(item.href);
              const Icon = item.icon;

              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={cn(
                    'flex items-center gap-3 px-4 py-3 rounded-lg transition-colors',
                    isActive
                      ? 'bg-strawberry-50 text-strawberry-700 font-medium'
                      : 'text-neutral-600 hover:bg-neutral-50'
                  )}
                >
                  <Icon className="w-5 h-5" />
                  <span>{item.name}</span>
                </Link>
              );
            })}
          </nav>
        </aside>

        {/* Main Content */}
        <main className="flex-1 overflow-hidden">
          {children}
        </main>
      </div>
    </div>
  );
}
