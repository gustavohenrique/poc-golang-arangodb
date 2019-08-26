import random
import json
import os

from faker import Faker
from faker.providers import file, address


fake = Faker()
fake.add_provider(file)
fake.add_provider(address)

wordlist1 = ['The best of', 'I like', 'I am sorry about', 'The most played', 'My favorite', 'I am tired of', 'Weekend', 'Friday', 'Play', 'Awesome', 'Strong', 'Fucking', 'Another', 'A few', 'Hot', 'Old', 'Traditional', 'Electrical', 'Interesting', 'Cute', 'Nice', 'Wonderful', 'Rare', 'Typical', 'Global', 'Dangerous', 'Powerful', 'Unusual', 'Pure', 'Aggressive', 'Boring', 'Massive', 'Angry', 'Local', 'Personal', 'Alternative', 'True', 'Positive', 'Negative', 'Red', 'Black', 'Green', 'Blue', 'Dark', 'Regular', 'Wrong', 'Extra', 'Unique', 'Classic', 'Private', 'Bright', 'Perfect', 'Correct', 'Slow', 'Fast', 'Fresh', 'Deep', 'Cool', 'Extreme', 'Exact', 'Lost', 'Sensitive', 'Weird', 'Dead', 'Wild', 'Adult', 'Sad', 'Strange', 'Sick', 'Crazy', 'Illegal', 'Funny', 'Royal', 'Sweet', 'Brave', 'Calm', 'Dirty', 'Honest', 'Brilliant', 'Drunk', 'Smart']
wordlist2 = ['Math', 'Algebra', 'Geometry', 'Science', 'Biology', 'Physics', 'Chemistry', 'Geography', 'History', 'Citizenship', 'Business', 'Home Economics', 'Art', 'Music', 'Politics', 'Technology']
wordlist3 = [i for i in range(1950, 2020)]


def _people(total):
    people = []
    unique = set()
    for i in range(1, total + 1):
        name = fake.name()
        if not name in unique:
            person = {
                '_key': '%d' % i,
                'name': name,
                'city': fake.city(),
            }
            people.append(person)
            unique.add(name)
    return people


def _playlists(total):
    l = []
    unique = set()
    for i in range(1, total + 1):
        subject = random.choice(wordlist2)
        name = '{} {} {}'.format(random.choice(wordlist1), random.choice(wordlist2), subject)
        if not name in unique:
            l.append({
                '_key': '%s' % i,
                'name': name,
            })
            unique.add(name)
    return l


def _audios(total):
    l = []
    for i in range(1, total + 1):
        filename = fake.file_name(extension='mp3')
        title = fake.text()
        l.append({
            '_key': '{}'.format(i),
            'title': title,
            'media': filename
        })
    return l


def main():
    max_data = 10
    students = _people(max_data)
    teachers = _people(100)
    playlists = _playlists(max_data)
    audios = _audios(max_data)
    subjects = []
    count = 1
    for w in wordlist2:
        subjects.append({'_key': '%s' % count, 'name': w})
        count = count + 1

    with open('json/subjects_vertex_.json', 'w') as f:
        json.dump(subjects, f)

    with open('json/audios_vertex_.json', 'w') as f:
        json.dump(audios, f)

    with open('json/teachers_vertex_.json', 'w') as f:
        json.dump(teachers, f)

    with open('json/students_vertex_.json', 'w') as f:
        json.dump(students, f)

    with open('json/playlists_vertex_.json', 'w') as f:
        json.dump(playlists, f)

    # Playlists can be tagged by subject
    with open('json/tagged_edge_.json', 'w') as f:
        tags = []
        i = 1
        for p in playlists:
            number = random.randint(1, 5)
            unique = set()
            for _ in range(number):
                s = random.choice(subjects)
                if not s['_key'] in unique:
                    tags.append({
                        '_key': '{}'.format(i),
                        '_from': 'playlists/%s' % p['_key'],
                        '_to': 'subjects/%s' % s['_key']
                    })
                    unique.add(s['_key'])
                    i = i + 1
        json.dump(tags, f)

    # Teacher records Audio
    with open('json/records_edge_.json', 'w') as f:
        records = []
        i = 1
        for row in audios:
            rand = random.randint(1,2)
            for _ in range(rand):
                teacher_id = random.choice(teachers)['_key']
                records.append({
                    '_key': '{}'.format(i),
                    '_from': 'teachers/%s' % teacher_id,
                    '_to': 'audios/%s' % row['_key']
                })
                i = i + 1
        json.dump(records, f)

    # Playlist plays Audio
    with open('json/plays_edge_.json', 'w') as f:
        plays = []
        i = 1
        for row in playlists:
            max_audios_in_playlist = random.randint(1, 10)
            ids = set()
            for x in range(max_audios_in_playlist):
                ids.add(random.randint(1, len(audios) - 1))
            for id in ids:
                plays.append({
                    '_key': '{}'.format(i),
                    '_from': 'playlists/%s' % row['_key'],
                    '_to': 'audios/%s' % id
                })
                i = i + 1
        json.dump(plays, f)

    # Student listen Playlist
    with open('json/listen_edge_.json', 'w') as f:
        listen = []
        i = 1
        for row in students:
            max_playlists = random.randint(1, 10)
            ids = set()
            for x in range(max_playlists):
                ids.add(random.randint(1, len(playlists) - 1))
            for id in ids:
                listen.append({
                    '_key': '{}'.format(i),
                    '_from': 'students/%s' % row['_key'],
                    '_to': 'playlists/%s' % id
                })
                i = i + 1
        json.dump(listen, f)


    print('Done!')


if __name__ == '__main__':
    main()

