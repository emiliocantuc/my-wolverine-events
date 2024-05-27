import sqlite3, argparse
import random


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description = 'Computes user recs and saves to database')
    parser.add_argument('--n', type = int, help = 'No. of events', default = 15, required = False)
    parser.add_argument('--output', type = str, help = 'Output db file to save the recs', default = 'data/main.db', required = False)
    args = parser.parse_args()

    try:
        conn = sqlite3.connect(args.output)
        cursor = conn.cursor()

        # Get all user IDs
        cursor.execute("SELECT user_id FROM users;")
        user_ids = cursor.fetchall()
        user_ids = [user_id[0] for user_id in user_ids]

        # Get latest week number
        cursor.execute('SELECT MAX(nweek) FROM statistics;')
        nweek = cursor.fetchone()[0] or 0 

        # Get evets from latest week
        cursor.execute("SELECT * FROM events WHERE nweek = ?", (nweek,))
        event_ids = [i[0] for i in cursor.fetchall()]
        print(f'Got # of users = {len(user_ids)}, nweek = {nweek}, # of events = {len(event_ids)}')

        # Randomly recommend (just to test)
        for user_id in user_ids:
            rec_event_ids = [random.choice(event_ids) for _ in range(args.n)]
            print(user_id, rec_event_ids)

            # Insert recommendations
            for rec_event_id in rec_event_ids:
                
                cursor.execute("""
                INSERT INTO recommended_events (user_id, event_id, method, params)
                VALUES (?, ?, ?, ?)
                """, (user_id, rec_event_id, 'random', ''))

            # Commit the transaction
            conn.commit()
        

    except sqlite3.Error as e:
        print(f"Database error: {e}")
    
    finally:
        conn.close()