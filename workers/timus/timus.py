import threading, time, json, requests, os
from datetime import datetime, timedelta
from sqlalchemy import create_engine, Column, Integer, String
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from bs4 import BeautifulSoup
from urllib.parse import quote

accounts = {}
password = "P@$$w0rd"
encoded_password = quote(password, safe='')
engine = create_engine(f'postgresql://bob:{encoded_password}@localhost:5432/test')

# Создание фабрики сессий
Session = sessionmaker(bind=engine)

# Определение базового класса моделей
Base = declarative_base()


# Определение моделей таблиц
class UnverifiedSubmission(Base):
    __tablename__ = 'unverified_submissions'

    submission_id = Column(String, primary_key=True)
    external_submission_id = Column(String)
    testing_system = Column(String)
    judge_id = Column(String)


class SubmissionVerdict(Base):
    __tablename__ = 'submissions_verdicts'

    id = Column(Integer, primary_key=True)
    sender_login = Column(String)
    task_id = Column(String)
    testing_system = Column(String)
    code = Column(String)
    submission_time = Column(String)
    contest_id = Column(Integer)
    contest_task_id = Column(Integer)
    verdict = Column(String)
    execution_time = Column(String)
    memory_used = Column(String)
    test = Column(String)
    language = Column(String)
    submission_number = Column(String)


class Submission(Base):
    __tablename__ = 'submissions'

    id = Column(Integer, primary_key=True)
    sender_login = Column(String)
    task_id = Column(String)
    testing_system = Column(String)
    code = Column(String)
    submission_time = Column(String)
    contest_id = Column(Integer)
    contest_task_id = Column(Integer)
    language = Column(String)
    sverdict_id = Column(String)


def timus_submission_sender(judgeid, account_name, code, language, task_id):
    global accounts
    # This is a function, that sends solve to timus
    # and returns the id of the submission
    # then we can use that id to check the verdict

    d = {
        "FreePascal 2.6": "31",
        "Visual C 2019": "63",
        "Visual C++ 2019": "64",
        "Visual C 2019 x64": "65",
        "Visual C++ 2019 x64": "66",
        "GCC 9.2 x64": "67",
        "G++ 9.2 x64": "68",
        "Clang++ 10 x64": "69",
        "Java 1.8": "32",
        "Visual C# 2019": "61",
        "Python 3.8 x64": "57",
        "PyPy 3.8 x64": "71",
        "Go 1.14 x64": "58",
        "Ruby 1.9": "18",
        "Haskell 7.6": "19",
        "Scala 2.11": "33",
        "Rust 1.58 x64": "72",
        "Kotlin 1.4.0": "60",
    }

    # creating a session
    r = requests.session()

    # this is a data of a post request to timus server
    data = {
        'action': 'submit',
        'SpaceID': 1,
        'JudgeID': judgeid,
        'Language': d[language],
        'ProblemNum': task_id,
        'Source': code,
    }

    # sending submission to timus server with our data
    r.post('https://acm.timus.ru/submit.aspx')

    txt = r.post('https://acm.timus.ru/submit.aspx', data=data)
    soup = BeautifulSoup(txt.text, 'html.parser')

    for row in soup.select("table.status.status_nofilter tr"):
        id_element = row.select_one("td.id")
        coder_element = row.select_one("td.coder a")
        problem_element = row.select_one("td.problem a")

        if id_element and coder_element and problem_element:
            id_value = id_element.text.strip()
            coder_value = coder_element.text.strip()
            problem_value = problem_element.text.strip().split(".", 1)[0]
            if coder_value == account_name and problem_value == task_id:
                return id_value

    with open('a.html', 'w') as f:
        f.write(txt.text)


def submitter():
    global accounts

    delta = timedelta(seconds=11)
    result_time = datetime.now() - delta

    for k, other_info in accounts.items():
        if "time" not in other_info:
            other_info["time"] = result_time

    # Получаем сначала все решения которые еще не были отправлены
    session = Session()

    # Получение всех записей, где testing_system равно "timus"
    submissions = session.query(Submission).filter(Submission.testing_system == "timus").all()

    for submission in submissions:
        all_busy = True

        # Затем пытаемся отправить каждое решение
        for JudgeId, other_info in accounts.items():
            if other_info["time"] + delta < datetime.now():
                other_info["time"] = datetime.now()
                all_busy = False

                # Используем judge_id k, чтобы послать решение
                # Получаем id посылки на тимусе
                id = timus_submission_sender(judgeid=JudgeId,
                                             account_name=other_info["name"],
                                             code=submission.code,
                                             language=submission.language,
                                             task_id=submission.task_id)

                # добавляем значение в UnverifiedSubmission
                new_submission = UnverifiedSubmission(
                    submission_id=submission.sverdict_id,
                    external_submission_id=id,
                    testing_system="timus",
                    judge_id=JudgeId
                )

                session.add(new_submission)
                session.commit()

                # изменяем значение submission_number в таблице SubmissionVerdict
                verdict = session.query(SubmissionVerdict).filter(
                    SubmissionVerdict.id == submission.sverdict_id).first()

                if verdict:
                    # Обновляем значения полей
                    verdict.verdict = "Compiling"
                    verdict.submission_number = id

                    session.commit()

                # удаляем из таблицы Submission
                session.delete(submission)
                session.commit()

        # если нет свободных -> ждем 0.5 секунды, вдруг освободятся
        if all_busy:
            break

    session.close()
    # TODO смотреть, освободится ли хотя бы одна штука через 0.5 секунды
    threading.Timer(0.5, submitter).start()


def construct_url(id, count=50):
    return f'https://acm.timus.ru/status.aspx?space=1&count={count}&from={id}'


def checker():
    session = Session()
    submissions = session.query(UnverifiedSubmission).filter(UnverifiedSubmission.testing_system == 'timus').all()

    # сортировка была произведена для удобства и для корректной работы алгоритма
    submissions.sort(key=lambda x: x.external_submission_id, reverse=True)
    submissions_dict = {submission.external_submission_id: {
        'submission_id': submission.submission_id,
        'testing_system': submission.testing_system,
        'judge_id': submission.judge_id
    } for submission in submissions}
    submissions_ids = [submission.external_submission_id for submission in submissions]

    while len(submissions) != 0:
        table = requests.get(construct_url(submissions[0].external_submission_id))

        soup = BeautifulSoup(table.text, 'html.parser')
        for link in soup.find_all('tr'):
            if link.get('class') != None and link.get('class')[0] in ['even', 'odd']:
                # получаем id посылки, которую мы щас парсим
                s_id = link.contents[0].text

                if s_id in submissions_ids:
                    id_in_other_table = submissions_dict[s_id]["submission_id"]
                    # удалить из словаря, списка id
                    del submissions_dict[s_id]
                    submissions_ids.remove(s_id)

                    # надо обновить всю информацию в таблице Verd....
                    submission = session.query(SubmissionVerdict) \
                        .filter_by(id=id_in_other_table).first()

                    if submission:
                        submission.verdict = str(link.contents[5].text)
                        submission.test = str(link.contents[6].text)
                        submission.execution_time = str(link.contents[7].text)
                        submission.memory_used = str(link.contents[8].text)
                        session.commit()

                        if link.contents[5].text in ['Compilation error', 'Wrong answer',
                                                     'Accepted', 'Time limit exceeded', 'Memory limit exceeded',
                                                     'Runtime error (non-zero exit code)', 'Runtime error']:
                            unver_sub = session.query(UnverifiedSubmission) \
                                .filter_by(submission_id=id_in_other_table).first()
                            if unver_sub:
                                session.delete(unver_sub)
                                session.commit()

        if len(submissions_ids) == 0:
            break

    session.close()

    threading.Timer(1, checker).start()


def main():
    global accounts
    
    with open(os.path.join(os.path.dirname(__file__), 'accounts.json'), 'r') as file:
        accounts = json.load(file)

    time.sleep(2)

    # Создаем объекты-потоки для каждой функции
    submitter_thread = threading.Thread(target=submitter)
    checker_thread = threading.Thread(target=checker)

    # Запускаем потоки
    submitter_thread.start()
    checker_thread.start()

    # Ожидаем завершения потоков
    submitter_thread.join()
    checker_thread.join()


if __name__ == "__main__":
    main()
