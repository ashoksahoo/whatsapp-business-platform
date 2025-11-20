import { useState, useEffect, useRef } from 'react';
import { Search, Send } from 'lucide-react';
import ContactCard from '../components/ContactCard';
import MessageBubble from '../components/MessageBubble';
import Avatar from '../components/Avatar';
import Input from '../components/Input';
import Button from '../components/Button';
import LoadingSpinner from '../components/LoadingSpinner';
import { useToast } from '../hooks/useToast';
import { apiClient } from '../services/api';
import { formatPhoneNumber } from '../lib/utils';
import type { Contact, Message } from '../types';

export default function MessagesPage() {
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [selectedContact, setSelectedContact] = useState<Contact | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [messageInput, setMessageInput] = useState('');
  const [isLoadingContacts, setIsLoadingContacts] = useState(true);
  const [isLoadingMessages, setIsLoadingMessages] = useState(false);
  const [isSending, setIsSending] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const toast = useToast();

  // Load contacts
  useEffect(() => {
    loadContacts();
  }, []);

  // Load messages when contact is selected
  useEffect(() => {
    if (selectedContact) {
      loadMessages(selectedContact.phone_number);
    }
  }, [selectedContact]);

  // Scroll to bottom when messages change
  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const loadContacts = async () => {
    try {
      setIsLoadingContacts(true);
      const response = await apiClient.listContacts({ limit: 50 });
      setContacts(response.data);
      if (response.data.length > 0 && !selectedContact) {
        setSelectedContact(response.data[0]);
      }
    } catch (error) {
      toast.error('Failed to load contacts');
      console.error('Error loading contacts:', error);
    } finally {
      setIsLoadingContacts(false);
    }
  };

  const loadMessages = async (phone: string) => {
    try {
      setIsLoadingMessages(true);
      const response = await apiClient.listMessages({ phone, limit: 100 });
      setMessages(response.data);
    } catch (error) {
      toast.error('Failed to load messages');
      console.error('Error loading messages:', error);
    } finally {
      setIsLoadingMessages(false);
    }
  };

  const handleSendMessage = async () => {
    if (!selectedContact || !messageInput.trim()) return;

    try {
      setIsSending(true);
      const message = await apiClient.sendMessage({
        to: selectedContact.phone_number,
        message_type: 'text',
        content: messageInput.trim(),
      });
      setMessages((prev) => [...prev, message]);
      setMessageInput('');
      toast.success('Message sent');
    } catch (error) {
      toast.error('Failed to send message');
      console.error('Error sending message:', error);
    } finally {
      setIsSending(false);
    }
  };

  const filteredContacts = contacts.filter(
    (contact) =>
      contact.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      contact.phone_number.includes(searchQuery)
  );

  return (
    <div className="h-full flex">
      {/* Conversation List */}
      <div className="w-[30%] border-r border-neutral-200 bg-white flex flex-col">
        <div className="p-4 border-b">
          <Input
            type="text"
            placeholder="Search contacts..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            icon={<Search className="w-5 h-5" />}
          />
        </div>
        <div className="flex-1 overflow-y-auto">
          {isLoadingContacts ? (
            <div className="flex justify-center items-center h-32">
              <LoadingSpinner />
            </div>
          ) : filteredContacts.length === 0 ? (
            <div className="text-center text-neutral-500 py-8">No contacts found</div>
          ) : (
            filteredContacts.map((contact) => (
              <ContactCard
                key={contact.id}
                contact={contact}
                active={selectedContact?.id === contact.id}
                onClick={() => setSelectedContact(contact)}
              />
            ))
          )}
        </div>
      </div>

      {/* Message Thread */}
      <div className="flex-1 flex flex-col bg-neutral-50">
        {selectedContact ? (
          <>
            {/* Header */}
            <div className="h-16 bg-white border-b px-6 flex items-center justify-between flex-shrink-0">
              <div className="flex items-center gap-3">
                <Avatar name={selectedContact.name} size="md" />
                <div>
                  <h2 className="font-semibold text-neutral-900">{selectedContact.name}</h2>
                  <p className="text-sm text-neutral-500">{formatPhoneNumber(selectedContact.phone_number)}</p>
                </div>
              </div>
            </div>

            {/* Messages */}
            <div className="flex-1 overflow-y-auto p-6">
              {isLoadingMessages ? (
                <div className="flex justify-center items-center h-full">
                  <LoadingSpinner />
                </div>
              ) : messages.length === 0 ? (
                <div className="text-center text-neutral-500">No messages yet</div>
              ) : (
                <>
                  {messages.map((message) => (
                    <MessageBubble key={message.id} message={message} />
                  ))}
                  <div ref={messagesEndRef} />
                </>
              )}
            </div>

            {/* Input */}
            <div className="bg-white border-t p-4 flex-shrink-0">
              <div className="flex items-end gap-3">
                <textarea
                  rows={1}
                  placeholder="Type a message..."
                  value={messageInput}
                  onChange={(e) => setMessageInput(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter' && !e.shiftKey) {
                      e.preventDefault();
                      handleSendMessage();
                    }
                  }}
                  className="flex-1 px-3 py-2 border border-neutral-300 rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-strawberry-500"
                />
                <Button
                  onClick={handleSendMessage}
                  disabled={!messageInput.trim()}
                  isLoading={isSending}
                  className="p-3"
                >
                  <Send className="w-5 h-5" />
                </Button>
              </div>
            </div>
          </>
        ) : (
          <div className="flex items-center justify-center h-full text-neutral-500">
            Select a contact to start messaging
          </div>
        )}
      </div>
    </div>
  );
}
