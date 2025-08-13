document.addEventListener('DOMContentLoaded', function() {
    const logoutBtn = document.getElementById('logoutBtn');
    const noteForm = document.getElementById('noteForm');
    const errorMessageEl = document.getElementById('errorMessage');
    const notesContainer = document.getElementById('notesContainer');

    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = '/';
        return;
    }

   
    loadNotes();

    logoutBtn.addEventListener('click', function() {
        localStorage.removeItem('token');
        window.location.href = '/';
    });

    noteForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const category = document.getElementById('category').value;
        const title = document.getElementById('title').value;
        const content = document.getElementById('content').value;

        try {
            const response = await fetch('/note', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    category: category,
                    title: title,
                    content: content
                })
            });

            const data = await response.json();

            if (response.ok) {
                
                showNotification('Note created successfully!', 'success');
                noteForm.reset();
                loadNotes();
            } else {
                errorMessageEl.textContent = data.error || 'Failed to create note';
            }
        } catch (err) {
            errorMessageEl.textContent = 'Error connecting to server';
            console.error('Note creation error:', err);
        }
    });

   
    async function loadNotes() {
        try {
            notesContainer.innerHTML = '<div class="loading">Loading notes...</div>';
            
            const response = await fetch('/note', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                throw new Error('Failed to fetch notes');
            }

            const notes = await response.json();
            renderNotes(notes);
        } catch (err) {
            errorMessageEl.textContent = 'Error loading notes';
            console.error('Error loading notes:', err);
            notesContainer.innerHTML = `<div class="error">Error loading notes: ${err.message}</div>`;
        }
    }

   
    function renderNotes(notes) {
        if (!notesContainer) return;
        
        if (!notes || notes.length === 0) {
            notesContainer.innerHTML = '<div class="empty-state">No notes yet. Create your first note!</div>';
            return;
        }

        notesContainer.innerHTML = '';
        notes.forEach(note => {
            const noteElement = document.createElement('div');
            noteElement.className = 'note-card';
            noteElement.innerHTML = `
                <span class="note-category">${note.category || 'General'}</span>
                <h3 class="note-title"><a href="/note/${note.id}/view?token=${token}">${note.title || 'Untitled Note'}</a></h3>
                <p class="note-content">${note.content || 'No content'}</p>
                <p class="note-date">Created: ${new Date(note.created_at).toLocaleString()}</p>
            `;
            notesContainer.appendChild(noteElement);
        });
    }

    
    function showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        
        document.body.appendChild(notification);
        
        setTimeout(() => {
            notification.classList.add('fade-out');
            setTimeout(() => notification.remove(), 500);
        }, 3000);
    }
});