import { useState, useEffect } from 'react';
import { Search, Edit } from 'lucide-react';
import Card from '../components/Card';
import Avatar from '../components/Avatar';
import Input from '../components/Input';
import LoadingSpinner from '../components/LoadingSpinner';
import { useToast } from '../hooks/useToast';
import { apiClient } from '../services/api';
import { formatPhoneNumber, formatTimestamp } from '../lib/utils';
import type { Contact } from '../types';

export default function ContactsPage() {
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const toast = useToast();

  useEffect(() => {
    loadContacts();
  }, []);

  const loadContacts = async () => {
    try {
      setIsLoading(true);
      const response = await apiClient.listContacts({ limit: 100 });
      setContacts(response.data);
    } catch (error) {
      toast.error('Failed to load contacts');
      console.error('Error loading contacts:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const filteredContacts = contacts.filter(
    (contact) =>
      contact.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      contact.phone_number.includes(searchQuery)
  );

  return (
    <div className="h-full flex flex-col">
      <div className="bg-white border-b px-6 py-4">
        <h1 className="text-2xl font-bold text-neutral-900 mb-4">Contacts</h1>
        <Input
          type="text"
          placeholder="Search contacts..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          icon={<Search className="w-5 h-5" />}
        />
      </div>

      <div className="flex-1 overflow-y-auto p-6">
        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <LoadingSpinner size="lg" />
          </div>
        ) : filteredContacts.length === 0 ? (
          <div className="text-center text-neutral-500 py-12">No contacts found</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {filteredContacts.map((contact) => (
              <Card key={contact.id} hover>
                <div className="flex items-start gap-4">
                  <Avatar name={contact.name} size="lg" />
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-neutral-900 truncate">{contact.name}</h3>
                    <p className="text-sm text-neutral-600">{formatPhoneNumber(contact.phone_number)}</p>
                    <div className="mt-3 space-y-1 text-sm text-neutral-500">
                      <div className="flex justify-between">
                        <span>Messages:</span>
                        <span className="font-medium">{contact.message_count}</span>
                      </div>
                      <div className="flex justify-between">
                        <span>Unread:</span>
                        <span className="font-medium text-strawberry-600">{contact.unread_count}</span>
                      </div>
                      {contact.last_message_at && (
                        <div className="flex justify-between">
                          <span>Last message:</span>
                          <span className="font-medium">{formatTimestamp(contact.last_message_at)}</span>
                        </div>
                      )}
                    </div>
                  </div>
                  <button className="text-neutral-400 hover:text-strawberry-600">
                    <Edit className="w-4 h-4" />
                  </button>
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
