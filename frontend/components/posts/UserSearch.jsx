'use client';

import { useState, useEffect, useRef } from 'react';
import { Search, X, User } from 'lucide-react';
import Image from 'next/image';

const UserSearch = ({ selectedUsers, onUserSelect, onUserRemove }) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [isSearching, setIsSearching] = useState(false);
  const [showDropdown, setShowDropdown] = useState(false);
  const searchRef = useRef(null);
  const dropdownRef = useRef(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (
        searchRef.current &&
        !searchRef.current.contains(event.target) &&
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target)
      ) {
        setShowDropdown(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Search for users with debouncing
  useEffect(() => {
    const searchUsers = async () => {
      if (searchQuery.trim().length < 2) {
        setSearchResults([]);
        setShowDropdown(false);
        return;
      }

      setIsSearching(true);
      try {
        const response = await fetch(`http://localhost:9000/users/search?q=${encodeURIComponent(searchQuery)}`, {
          credentials: 'include',
        });

        if (response.ok) {
          const users = await response.json();
          // Filter out already selected users
          const filteredUsers = users.filter(
            user => !selectedUsers.some(selected => selected.id === user.id)
          );
          setSearchResults(filteredUsers);
          setShowDropdown(filteredUsers.length > 0);
        } else {
          setSearchResults([]);
          setShowDropdown(false);
        }
      } catch (error) {
        console.error('Error searching users:', error);
        setSearchResults([]);
        setShowDropdown(false);
      } finally {
        setIsSearching(false);
      }
    };

    const timeoutId = setTimeout(searchUsers, 300);
    return () => clearTimeout(timeoutId);
  }, [searchQuery, selectedUsers]);

  const handleUserSelect = (user) => {
    onUserSelect(user);
    setSearchQuery('');
    setSearchResults([]);
    setShowDropdown(false);
  };

  const getDisplayName = (user) => {
    if (user.nickname) {
      return user.nickname;
    }
    if (user.firstname && user.lastname) {
      return `${user.firstname} ${user.lastname}`;
    }
    return user.firstname || user.lastname || 'Unknown User';
  };

  const getFullName = (user) => {
    if (user.firstname && user.lastname) {
      return `${user.firstname} ${user.lastname}`;
    }
    return user.firstname || user.lastname || '';
  };

  return (
    <div className="space-y-3">
      {/* Selected Users */}
      {selectedUsers.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {selectedUsers.map((user) => (
            <div
              key={user.id}
              className="flex items-center gap-2 px-3 py-1 rounded-full text-sm"
              style={{ backgroundColor: 'var(--tertiary-background)', color: 'var(--quaternary-text)' }}
            >
              <div className="w-6 h-6 rounded-full flex items-center justify-center" style={{ backgroundColor: 'var(--tertiary-background)' }}>
                {user.avatar && user.avatar !== "no profile photo" ? (
                  <Image
                    src={`http://localhost:9000/avatar?avatar=${encodeURIComponent(user.avatar)}`}
                    alt={getDisplayName(user)}
                    width={24}
                    height={24}
                    className="w-6 h-6 rounded-full object-cover"
                  />
                ) : (
                  <User className="w-4 h-4" style={{ color: 'var(--tertiary-text)' }} />
                )}
              </div>
              <span>{getDisplayName(user)}</span>
              <button
                type="button"
                onClick={() => onUserRemove(user.id)}
                style={{ color: 'var(--tertiary-text)' }}
                onMouseOver={(e) => e.currentTarget.style.color = 'var(--quinary-text)'}
                onMouseOut={(e) => e.currentTarget.style.color = 'var(--tertiary-text)'}
              >
                <X className="w-4 h-4" />
              </button>
            </div>
          ))}
        </div>
      )}

      {/* Search Input */}
      <div className="relative" ref={searchRef}>
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4" style={{ color: 'var(--secondary-text)' }} />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search for people to share with..."
            className="w-full pl-10 pr-4 py-2 border rounded-lg focus:ring-2 focus:border-transparent"
            style={{ borderColor: 'var(--border-color)', '--ring-color': 'var(--primary-accent)' }}
            onFocus={(e) => {
              e.target.style.borderColor = 'var(--primary-accent)';
              if (searchResults.length > 0) {
                setShowDropdown(true);
              }
            }}
            onBlur={(e) => e.target.style.borderColor = 'var(--border-color)'}
          />
          {isSearching && (
            <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2" style={{ borderColor: 'var(--primary-accent)' }}></div>
            </div>
          )}
        </div>

        {/* Search Results Dropdown */}
        {showDropdown && (
          <div
            ref={dropdownRef}
            className="absolute z-10 w-full mt-1 rounded-lg shadow-lg max-h-60 overflow-y-auto"
            style={{ backgroundColor: 'var(--secondary-background)', border: '1px solid var(--border-color)' }}
          >
            {searchResults.map((user) => (
              <button
                key={user.id}
                type="button"
                onClick={() => handleUserSelect(user)}
                className="w-full px-4 py-3 text-left flex items-center gap-3 border-b last:border-b-0"
                style={{ borderColor: 'var(--border-color)' }}
                onMouseOver={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-background)'}
                onMouseOut={(e) => e.currentTarget.style.backgroundColor = 'var(--secondary-background)'}
              >
                <div className="w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0" style={{ backgroundColor: 'var(--tertiary-background)' }}>
                  {user.avatar && user.avatar !== "no profile photo" ? (
                    <Image
                      src={`http://localhost:9000/avatar?avatar=${encodeURIComponent(user.avatar)}`}
                      alt={getDisplayName(user)}
                      width={32}
                      height={32}
                      className="w-8 h-8 rounded-full object-cover"
                    />
                  ) : (
                    <User className="w-5 h-5" style={{ color: 'var(--quaternary-text)' }} />
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="font-medium truncate" style={{ color: 'var(--primary-text)' }}>
                    {getDisplayName(user)}
                  </div>
                  {user.nickname && getFullName(user) && (
                    <div className="text-sm truncate" style={{ color: 'var(--secondary-text)' }}>
                      {getFullName(user)}
                    </div>
                  )}
                </div>
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default UserSearch;
