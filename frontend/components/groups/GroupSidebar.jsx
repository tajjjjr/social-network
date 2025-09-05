'use client';

import { useState, useEffect } from 'react';
import { Calendar, Users, Settings } from 'lucide-react';

export default function GroupSidebar({ group }) {
  const [events, setEvents] = useState([]);
  const [members, setMembers] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGroupData = async () => {
      try {
        // Fetch group events
        const eventsResponse = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${group.id}/events`, {
          credentials: 'include'
        });
        if (eventsResponse.ok) {
          const eventsData = await eventsResponse.json();
          setEvents(eventsData || []);
        }

        // Fetch group members
        const membersResponse = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/groups/${group.id}/members`, {
          credentials: 'include'
        });
        if (membersResponse.ok) {
          const membersData = await membersResponse.json();
          setMembers(membersData || []);
        }
      } catch (error) {
        console.error('Error fetching group data:', error);
      } finally {
        setLoading(false);
      }
    };

    if (group?.id) {
      fetchGroupData();
    }
  }, [group?.id]);

  return (
    <div className="space-y-6">
      {/* Group Info Card */}
      <div 
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--secondary-background)' }}
      >
        <h3 className="text-lg font-bold mb-4" style={{ color: 'var(--primary-text)' }}>
          About
        </h3>
        <div className="space-y-3">
          <div className="flex items-center gap-3">
            <Users className="w-5 h-5" style={{ color: 'var(--secondary-text)' }} />
            <span style={{ color: 'var(--secondary-text)' }}>
              {members.length} members
            </span>
          </div>
          <div className="flex items-center gap-3">
            <Calendar className="w-5 h-5" style={{ color: 'var(--secondary-text)' }} />
            <span style={{ color: 'var(--secondary-text)' }}>
              Created {new Date(group.created_at).toLocaleDateString()}
            </span>
          </div>
          <div className="flex items-center gap-3">
            <Settings className="w-5 h-5" style={{ color: 'var(--secondary-text)' }} />
            <span style={{ color: 'var(--secondary-text)' }}>
              {group.privacy} group
            </span>
          </div>
        </div>
      </div>

      {/* Upcoming Events Card */}
      <div 
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--secondary-background)' }}
      >
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-bold" style={{ color: 'var(--primary-text)' }}>
            Upcoming Events
          </h3>
          <button
            className="text-sm"
            style={{ color: 'var(--primary-accent)' }}
          >
            See all
          </button>
        </div>
        
        {loading ? (
          <div className="text-center py-4">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 mx-auto" style={{ borderColor: 'var(--primary-accent)' }}></div>
          </div>
        ) : events.length > 0 ? (
          <div className="space-y-3">
            {events.slice(0, 3).map(event => (
              <div key={event.id} className="flex gap-3 p-3 rounded-lg hover:bg-opacity-50" style={{ backgroundColor: 'var(--primary-background)' }}>
                <div className="flex-shrink-0">
                  <div className="w-12 h-12 rounded-lg flex flex-col items-center justify-center text-xs" style={{ backgroundColor: 'var(--primary-accent)', color: 'white' }}>
                    <span className="font-bold">{new Date(event.event_time).getDate()}</span>
                    <span>{new Date(event.event_time).toLocaleDateString('en', { month: 'short' })}</span>
                  </div>
                </div>
                <div className="flex-1 min-w-0">
                  <h4 className="font-medium truncate" style={{ color: 'var(--primary-text)' }}>
                    {event.title}
                  </h4>
                  <p className="text-sm truncate" style={{ color: 'var(--secondary-text)' }}>
                    {new Date(event.event_time).toLocaleTimeString('en', { hour: '2-digit', minute: '2-digit' })}
                  </p>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <Calendar className="w-12 h-12 mx-auto mb-3 opacity-50" style={{ color: 'var(--secondary-text)' }} />
            <p className="text-sm" style={{ color: 'var(--secondary-text)' }}>
              No upcoming events
            </p>
          </div>
        )}
      </div>

      {/* Recent Members Card */}
      <div 
        className="rounded-xl p-6"
        style={{ backgroundColor: 'var(--secondary-background)' }}
      >
        <div className="flex justify-between items-center mb-4">
          <h3 className="text-lg font-bold" style={{ color: 'var(--primary-text)' }}>
            Members
          </h3>
          <button
            className="text-sm"
            style={{ color: 'var(--primary-accent)' }}
          >
            See all
          </button>
        </div>
        
        {loading ? (
          <div className="text-center py-4">
            <div className="animate-spin rounded-full h-6 w-6 border-b-2 mx-auto" style={{ borderColor: 'var(--primary-accent)' }}></div>
          </div>
        ) : members.length > 0 ? (
          <div className="space-y-3">
            {members.slice(0, 5).map(member => (
              <div key={member.id} className="flex items-center gap-3">
                <img
                  src={member.avatar || '/default-avatar.png'}
                  alt={member.name}
                  className="w-10 h-10 rounded-full object-cover"
                />
                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate" style={{ color: 'var(--primary-text)' }}>
                    {member.firstname} {member.lastname}
                  </p>
                  <p className="text-sm truncate" style={{ color: 'var(--secondary-text)' }}>
                    @{member.nickname}
                  </p>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <Users className="w-12 h-12 mx-auto mb-3 opacity-50" style={{ color: 'var(--secondary-text)' }} />
            <p className="text-sm" style={{ color: 'var(--secondary-text)' }}>
              No members yet
            </p>
          </div>
        )}
      </div>
    </div>
  );
}